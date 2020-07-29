package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	remote_file_storage "github.com/jochenboesmans/gedcom-parser/remote-file-storage"
	"github.com/jochenboesmans/gedcom-parser/server"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"github.com/jochenboesmans/gedcom-parser/model"
	"github.com/jochenboesmans/gedcom-parser/util"
)

var from = flag.String("from", "ged", "type of file to parse")
var to = flag.String("to", "json", "type of file to create")

func main() {
	switch os.Args[1] {
	case "parse":
		parse()
	case "serve":
		server.Serve()
	}
}

func parse() {
	flag.Parse()
	beginTime := time.Now()

	files, err := ioutil.ReadDir("io")
	if err != nil {
		log.Fatal("Unable to read from folder ./io")
	}

	var concurrentlyOpenFiles = make(chan int, 1020)
	waitGroup := &sync.WaitGroup{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), fmt.Sprintf(".%s", *from)) {
			waitGroup.Add(1)
			switch *from {
			case "ged":
				go parseGedcom(f.Name(), waitGroup, concurrentlyOpenFiles)
			case "json", "protobuf":
				go parseJsonOrProtobuf(f.Name(), waitGroup, concurrentlyOpenFiles)
			}
		}
	}
	waitGroup.Wait()

	fmt.Printf("total time taken: %f second.\n", float64(time.Since(beginTime))*math.Pow10(-9))
}

func parseJsonOrProtobuf(inputFileName string, outerWaitGroup *sync.WaitGroup, concurrentlyOpenFiles chan int) {
	gedcom := &model.Gedcom{}
	switch *from {
	case "json":
		readJSON(gedcom, concurrentlyOpenFiles, inputFileName)
	case "protobuf":
		readProtobuf(gedcom, concurrentlyOpenFiles, inputFileName)
	}

	//writeToGedcom(gedcom, concurrentlyOpenFiles, inputFileName)
	writeToGedcomS3(gedcom, inputFileName)

	outerWaitGroup.Done()
}

func readJSON(gedcom *model.Gedcom, concurrentlyOpenFiles chan int, inputFileName string) {
	concurrentlyOpenFiles <- 1 // premature increment of semaphore to prevent race condition
	jsonFile, err := ioutil.ReadFile("./io/" + inputFileName)
	if err != nil {
		<-concurrentlyOpenFiles
		log.Print(err)
	}

	err = json.Unmarshal(jsonFile, gedcom)
	util.Check(err)
}

func readProtobuf(gedcom *model.Gedcom, concurrentlyOpenFiles chan int, inputFileName string) {
	concurrentlyOpenFiles <- 1 // premature increment of semaphore to prevent race condition
	protobufFile, err := ioutil.ReadFile("./io/" + inputFileName)
	if err != nil {
		<-concurrentlyOpenFiles
		log.Print(err)
	}

	err = proto.Unmarshal(protobufFile, gedcom)
	util.Check(err)
}

func writeToGedcomS3(gedcom *model.Gedcom, inputFileName string) {
	_, err := remote_file_storage.S3Write(inputFileName, writableGedcom(gedcom))
	util.Check(err)
}

func writableGedcom(gedcom *model.Gedcom) *[]byte {
	gedcomString := ""

	header := "0 HEAD\n"
	gedcomString += header

	for _, i := range gedcom.Individuals {
		firstLine := fmt.Sprintf("0 %s INDI\n", i.Id)
		gedcomString += firstLine

		for _, n := range i.Names {
			nameLine := fmt.Sprintf("1 NAME %s/%s/\n", n.GivenName, n.Surname)
			gedcomString += nameLine
		}

		genderMap := map[string]string{
			"MALE":   "M",
			"FEMALE": "F",
		}
		genderLine := fmt.Sprintf("1 SEX %s\n", genderMap[i.Gender])
		gedcomString += genderLine
	}

	for _, f := range gedcom.Families {
		firstLine := fmt.Sprintf("0 %s FAM\n", f.Id)
		gedcomString += firstLine

		if f.FatherId != "" {
			fatherLine := fmt.Sprintf("1 HUSB %s\n", f.FatherId)
			gedcomString += fatherLine
		}
		if f.MotherId != "" {
			motherLine := fmt.Sprintf("1 WIFE %s\n", f.MotherId)
			gedcomString += motherLine
		}

		for _, childId := range f.ChildIds {
			childLine := fmt.Sprintf("1 CHIL %s\n", childId)
			gedcomString += childLine
		}
	}

	trailer := "0 TRLR\n"
	gedcomString += trailer

	gedcomBytes := []byte(gedcomString)

	return &gedcomBytes

}

func writeToGedcom(gedcom *model.Gedcom, concurrentlyOpenFiles chan int, inputFileName string) {
	concurrentlyOpenFiles <- 1
	writeFile, err := os.Create("./io/generated-" + strings.Split(inputFileName, ".")[0] + ".ged")
	if err != nil {
		<-concurrentlyOpenFiles
		log.Print(err)
	}

	w := bufio.NewWriter(writeFile)

	gedcomString := writableGedcom(gedcom)

	_, err = w.Write(*gedcomString)
	util.Check(err)

	err = writeFile.Close()
	if err != nil {
		log.Print(err)
	} else {
		<-concurrentlyOpenFiles
	}

}

func parseGedcom(inputFileName string, outerWaitGroup *sync.WaitGroup, concurrentlyOpenFiles chan int) {
	gedcom := readGedcom(concurrentlyOpenFiles, inputFileName)

	switch *to {
	case "json":
		writeGedcomToJSON(gedcom, concurrentlyOpenFiles, inputFileName)
	case "protobuf":
		writeGedcomToProtobuf(gedcom, concurrentlyOpenFiles, inputFileName)
	}

	outerWaitGroup.Done()
}

func readGedcom(concurrentlyOpenFiles chan int, inputFileName string) *model.ConcurrencySafeGedcom {
	concurrentlyOpenFiles <- 1 // premature increment of semaphore to prevent race condition
	file, err := os.Open("./io/" + inputFileName)
	if err != nil {
		<-concurrentlyOpenFiles
		log.Print(err)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	recordLines := []*gedcomSpec.Line{}
	waitGroup := &sync.WaitGroup{}

	gedcom := model.NewConcurrencySafeGedcom()

	i := 0
	for fileScanner.Scan() {
		line := ""
		if i == 0 {
			line = strings.TrimPrefix(fileScanner.Text(), "\uFEFF")
		} else {
			line = fileScanner.Text()
		}
		gedcomLine := gedcomSpec.NewLine(&line)

		// interpret record once it's fully read
		if len(recordLines) > 0 && *gedcomLine.Level() == 0 {
			waitGroup.Add(1)
			go gedcom.InterpretRecord(recordLines, waitGroup)
			recordLines = []*gedcomSpec.Line{}
		}
		recordLines = append(recordLines, gedcomLine)
		i++
	}

	waitGroup.Wait()
	err = file.Close()
	if err != nil {
		log.Print(err)
	} else {
		<-concurrentlyOpenFiles
	}

	return gedcom
}

func writeGedcomToProtobuf(g *model.ConcurrencySafeGedcom, concurrentlyOpenFiles chan int, inputFileName string) {
	gedcomProtobuf, err := proto.Marshal(&g.Gedcom)
	util.Check(err)
	_, err = remote_file_storage.S3Write(inputFileName, &gedcomProtobuf)
	util.Check(err)

	//concurrentlyOpenFiles <- 1
	//protobufFile, err := os.Create("./io/generated-" + strings.Split(inputFileName, ".")[0] + ".protobuf")
	//if err != nil {
	//	<-concurrentlyOpenFiles
	//	log.Print(err)
	//}
	//protoWriter := bufio.NewWriter(protobufFile)
	//_, err = protoWriter.Write(gedcomProtobuf)
	//util.Check(err)
	//err = protoWriter.Flush()
	//util.Check(err)
	//
	//err = protobufFile.Close()
	//if err != nil {
	//	log.Print(err)
	//} else {
	//	<-concurrentlyOpenFiles
	//}
}

func writeGedcomToJSON(g *model.ConcurrencySafeGedcom, concurrentlyOpenFiles chan int, inputFileName string) {
	gedcomJson, err := json.Marshal(g.Gedcom)
	util.Check(err)
	_, err = remote_file_storage.S3Write(inputFileName, &gedcomJson)
	util.Check(err)

	//concurrentlyOpenFiles <- 1
	//jsonFile, err := os.Create("./io/generated-" + strings.Split(inputFileName, ".")[0] + ".json")
	//if err != nil {
	//	<-concurrentlyOpenFiles
	//	log.Print(err)
	//}
	//jsonWriter := bufio.NewWriter(jsonFile)
	//_, err = jsonWriter.Write(gedcomJson)
	//util.Check(err)
	//err = jsonWriter.Flush()
	//util.Check(err)
	//
	//err = jsonFile.Close()
	//if err != nil {
	//	log.Print(err)
	//} else {
	//	<-concurrentlyOpenFiles
	//}
}
