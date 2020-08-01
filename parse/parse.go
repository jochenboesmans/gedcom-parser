package parse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"io"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Parse(inputFilePath string, outputFilePath string) {
	beginTime := time.Now()

	input, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open input file at %s\n", inputFilePath)
	}

	var output *[]byte
	inputReader := bytes.NewReader(input)

	switch filepath.Ext(inputFilePath) {
	case ".ged":
		output, err = ParseGedcom(inputReader, outputFilePath)
	case ".json":
		output, err = ParseJSON(inputReader)
	case ".protobuf":
		output, err = ParseProtobuf(inputReader)
	default:
		log.Fatalf("failed to match input file extension to: .ged|.json|.protobuf\n")
	}

	err = ioutil.WriteFile(outputFilePath, *output, 0600)
	if err != nil {
		log.Fatalf("failed to write to output file at %s\n", outputFilePath)
	}

	fmt.Printf("total time taken: %f second.\n", float64(time.Since(beginTime))*math.Pow10(-9))
}

func ParseGedcom(inputReader io.Reader, to string) (*[]byte, error) {
	fileScanner := bufio.NewScanner(inputReader)
	fileScanner.Split(bufio.ScanLines)

	recordLines := []*gedcomSpec.Line{}
	waitGroup := &sync.WaitGroup{}

	gedcom := gedcomSpec.NewConcurrencySafeGedcom()

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

	switch filepath.Ext(to) {
	case ".json":
		return GedcomToJSON(gedcom)
	case ".protobuf":
		return GedcomToProto(gedcom)
	}

	return nil, fmt.Errorf("failed to match output file extension to: %s", ".json|.protobuf")
}

func ParseJSON(inputReader io.Reader) (*[]byte, error) {
	gedcom := &gedcomSpec.Gedcom{}

	gedcomJson, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(gedcomJson, gedcom)
	if err != nil {
		return nil, err
	}

	gedcomBuf := WritableGedcom(gedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}

func ParseProtobuf(inputReader io.Reader) (*[]byte, error) {
	var gedcom *gedcomSpec.Gedcom

	gedcomProto, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(gedcomProto, gedcom)
	if err != nil {
		return nil, err
	}

	gedcomBuf := WritableGedcom(gedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}

func WritableGedcom(gedcom *gedcomSpec.Gedcom) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})

	header := "0 HEAD\n"
	buf.WriteString(header)

	for _, i := range gedcom.Individuals {
		firstLine := fmt.Sprintf("0 %s INDI\n", i.Id)
		buf.WriteString(firstLine)

		for _, n := range i.Names {
			nameLine := fmt.Sprintf("1 NAME %s/%s/\n", n.GivenName, n.Surname)
			buf.WriteString(nameLine)
		}

		genderMap := map[string]string{
			"MALE":   "M",
			"FEMALE": "F",
		}
		genderLine := fmt.Sprintf("1 SEX %s\n", genderMap[i.Gender])
		buf.WriteString(genderLine)
	}

	for _, f := range gedcom.Families {
		firstLine := fmt.Sprintf("0 %s FAM\n", f.Id)
		buf.WriteString(firstLine)

		if f.FatherId != "" {
			fatherLine := fmt.Sprintf("1 HUSB %s\n", f.FatherId)
			buf.WriteString(fatherLine)
		}
		if f.MotherId != "" {
			motherLine := fmt.Sprintf("1 WIFE %s\n", f.MotherId)
			buf.WriteString(motherLine)
		}

		for _, childId := range f.ChildIds {
			childLine := fmt.Sprintf("1 CHIL %s\n", childId)
			buf.WriteString(childLine)
		}
	}

	trailer := "0 TRLR\n"
	buf.WriteString(trailer)

	return buf
}

func GedcomToJSON(gedcom *gedcomSpec.ConcurrencySafeGedcom) (*[]byte, error) {
	gedcomJson, err := json.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomJson, nil
}

func GedcomToProto(gedcom *gedcomSpec.ConcurrencySafeGedcom) (*[]byte, error) {
	gedcomProto, err := proto.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomProto, nil
}
