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

func trimBOM(line string) string {
	return strings.TrimPrefix(line, "\uFEFF")
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
		readLine := fileScanner.Text()
		if i == 0 {
			line = trimBOM(readLine)
		} else {
			line = readLine
		}
		gedcomLine := gedcomSpec.NewLine(&line)

		level, err := gedcomLine.Level()
		if err != nil {
			continue
		}

		// interpret record once it's fully read
		if len(recordLines) > 0 && level == 0 {
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
		return gedcom.ToJSON()
	case ".protobuf":
		return gedcom.ToProto()
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

	gedcomBuf := gedcomSpec.WritableGedcom(concSafeGedcom)
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

	gedcomBuf := gedcomSpec.WritableGedcom(concSafeGedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}
