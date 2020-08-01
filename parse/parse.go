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

	gedcom.RemoveInvalidFamilies()

	switch filepath.Ext(to) {
	case ".json":
		return GedcomToJSON(gedcom)
	case ".protobuf":
		return GedcomToProto(gedcom)
	}

	return nil, fmt.Errorf("failed to match output file extension to: %s", ".json|.protobuf")
}

func ParseJSON(inputReader io.Reader) (*[]byte, error) {
	gedcom := gedcomSpec.Gedcom{}

	gedcomJson, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(gedcomJson, &gedcom)
	if err != nil {
		return nil, err
	}

	concSafeGedcom := gedcomSpec.NewConcurrencySafeGedcom()
	concSafeGedcom.Gedcom = gedcom

	concSafeGedcom.RemoveInvalidFamilies()

	gedcomBuf := WritableGedcom(concSafeGedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}

func ParseProtobuf(inputReader io.Reader) (*[]byte, error) {
	var gedcom gedcomSpec.Gedcom

	gedcomProto, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(gedcomProto, &gedcom)
	if err != nil {
		return nil, err
	}

	concSafeGedcom := gedcomSpec.NewConcurrencySafeGedcom()
	concSafeGedcom.Gedcom = gedcom

	concSafeGedcom.RemoveInvalidFamilies()

	gedcomBuf := WritableGedcom(concSafeGedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}

func WritableGedcom(concSafeGedcom *gedcomSpec.ConcurrencySafeGedcom) *bytes.Buffer {
	gedcom := concSafeGedcom.Gedcom
	buf := bytes.NewBuffer([]byte{})

	header := "0 HEAD\n"
	buf.WriteString(header)

	for _, i := range gedcom.Individuals {
		firstLine := fmt.Sprintf("0 %s INDI\n", i.Id)
		buf.WriteString(firstLine)

		for _, n := range i.Names {
			nameLine := fmt.Sprintf("1 NAME %s/%s/\n", n.GivenName, n.Surname)
			buf.WriteString(nameLine)

			primaryLine := fmt.Sprintf("2 _PRIM %s\n", primaryValueByBool[n.Primary])
			buf.WriteString(primaryLine)
		}

		if i.BirthDate != nil {
			firstLine := fmt.Sprintf("1 BIRT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if i.BirthDate.Year != 0 && i.BirthDate.Month != 0 && i.BirthDate.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", i.BirthDate.Day, monthAbbrByInt[int(i.BirthDate.Month)], i.BirthDate.Year)
			} else if i.BirthDate.Year != 0 && i.BirthDate.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", monthAbbrByInt[int(i.BirthDate.Month)], i.BirthDate.Year)
			} else if i.BirthDate.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", i.BirthDate.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}
		}
		if i.DeathDate != nil {
			firstLine := fmt.Sprintf("1 DEAT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if i.DeathDate.Year != 0 && i.DeathDate.Month != 0 && i.BirthDate.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", i.DeathDate.Day, monthAbbrByInt[int(i.DeathDate.Month)], i.DeathDate.Year)
			} else if i.DeathDate.Year != 0 && i.DeathDate.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", monthAbbrByInt[int(i.DeathDate.Month)], i.BirthDate.Year)
			} else if i.DeathDate.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", i.DeathDate.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}
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

var monthAbbrByInt = map[int]string{
	1:  "JAN",
	2:  "FEB",
	3:  "MAR",
	4:  "APR",
	5:  "MAY",
	6:  "JUN",
	7:  "JUL",
	8:  "AUG",
	9:  "SEP",
	10: "OCT",
	11: "NOV",
	12: "DEC",
}

var primaryValueByBool = map[bool]string{
	true:  "Y",
	false: "N",
}
