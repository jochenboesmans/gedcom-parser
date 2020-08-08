package parse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"github.com/jochenboesmans/gedcom-parser/util"
	"io"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
Parse local files representing a gedcom structure to a different format representing the same structure.

Example usage: Parse("./familytree.ged", "./familytree.json") would parse the GEDCOM file at ./familytree.ged into a json structure and put the result in a file at ./familytree.json.
*/
func Parse(inputFilePath string, outputFilePath string) {
	beginTime := time.Now()

	input, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		log.Fatalf("failed to read from input file at %s with error: %s\n", inputFilePath, err)
	}

	var output *[]byte
	inputReader := bytes.NewReader(input)

	switch filepath.Ext(inputFilePath) {
	case ".ged":
		output, err = ParseGedcom(inputReader, outputFilePath)
		if err != nil {
			log.Fatalf("failed to parse GEDCOM file at %s with error: %s\n", inputFilePath, err)
		}
	case ".json":
		output, err = ParseJSON(inputReader)
		if err != nil {
			log.Fatalf("failed to parse JSON file at %s with error: %s\n", inputFilePath, err)
		}
	case ".protobuf":
		output, err = ParseProtobuf(inputReader)
		if err != nil {
			log.Fatalf("failed to parse Protobuf file at %s with error: %s\n", inputFilePath, err)
		}
	default:
		log.Fatalf("failed to match input file (at %s) extension to: .ged|.json|.protobuf\n", inputFilePath)
	}

	err = ioutil.WriteFile(outputFilePath, *output, 0600)
	if err != nil {
		log.Fatalf("failed to write to output file at %s with error: %s\n", outputFilePath, err)
	}

	secondsSinceBeginTime := float64(time.Since(beginTime)) * math.Pow10(-9)
	fmt.Printf("successfully parsed file at %s to %s. total time taken: %f seconds\n", inputFilePath, outputFilePath, secondsSinceBeginTime)
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

		level, err := gedcomLine.Level()
		// interpret record once it's fully read
		if len(recordLines) > 0 && err == nil && level == 0 {
			waitGroup.Add(1)
			go gedcom.InterpretRecord(recordLines, waitGroup)
			recordLines = []*gedcomSpec.Line{}
		}
		recordLines = append(recordLines, gedcomLine)
		i++
	}

	waitGroup.Wait()

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
	// try to decode non-utf8 fields, keep encoded version if it fails
	_ = concSafeGedcom.DecodeUnicodeFields()

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

			primaryLine := fmt.Sprintf("2 _PRIM %s\n", util.PrimaryValueByBool[n.Primary])
			buf.WriteString(primaryLine)
		}

		for _, b := range i.BirthEvents {
			firstLine := fmt.Sprintf("1 BIRT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if b.Date.Year != 0 && b.Date.Month != 0 && b.Date.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", b.Date.Day, util.MonthAbbrByInt[int(b.Date.Month)], b.Date.Year)
			} else if b.Date.Year != 0 && b.Date.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", util.MonthAbbrByInt[int(b.Date.Month)], b.Date.Year)
			} else if b.Date.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", b.Date.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}

			if b.Place != "" {
				placeLine := fmt.Sprintf("2 PLAC %s\n", b.Place)
				buf.WriteString(placeLine)
			}

			primaryLine := fmt.Sprintf("2 _PRIM %s\n", util.PrimaryValueByBool[b.Primary])
			buf.WriteString(primaryLine)
		}

		for _, d := range i.DeathEvents {
			firstLine := fmt.Sprintf("1 DEAT\n")
			buf.WriteString(firstLine)

			var secondLine string
			if d.Date.Year != 0 && d.Date.Month != 0 && d.Date.Day != 0 {
				secondLine = fmt.Sprintf("2 DATE %d %s %d\n", d.Date.Day, util.MonthAbbrByInt[int(d.Date.Month)], d.Date.Year)
			} else if d.Date.Year != 0 && d.Date.Month != 0 {
				secondLine = fmt.Sprintf("2 DATE %s %d\n", util.MonthAbbrByInt[int(d.Date.Month)], d.Date.Year)
			} else if d.Date.Year != 0 {
				secondLine = fmt.Sprintf("2 DATE %d\n", d.Date.Year)
			}
			if secondLine != "" {
				buf.WriteString(firstLine)
			}

			if d.Place != "" {
				placeLine := fmt.Sprintf("2 PLAC %s\n", d.Place)
				buf.WriteString(placeLine)
			}

			if primaryValue, hit := util.PrimaryValueByBool[d.Primary]; hit {
				primaryLine := fmt.Sprintf("2 _PRIM %s\n", primaryValue)
				buf.WriteString(primaryLine)
			}
		}

		if genderLetter, hit := util.GenderLetterByFull[i.Gender]; hit {
			genderLine := fmt.Sprintf("1 SEX %s\n", genderLetter)
			buf.WriteString(genderLine)
		}
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
