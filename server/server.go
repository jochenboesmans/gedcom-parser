package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"github.com/jochenboesmans/gedcom-parser/model"
	remote_file_storage "github.com/jochenboesmans/gedcom-parser/remote-file-storage"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

func Serve() {
	addr := ":8080"
	http.HandleFunc("/parse", handleParse)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("failed to start server at %s: %s", addr, err)
	}
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseParams(w, r)
	if err != nil {
		return
	}
	input, err := remote_file_storage.S3Read(from)
	if err != nil {
		_, _ = fmt.Fprintf(w, "failed to read from S3: %s", err)
		return
	}

	inputReader := bytes.NewReader(*input)

	switch filepath.Ext(from) {
	case ".ged":
		output, err := parseGedcom(inputReader, to)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to parse gedcom: %s", err)
			return
		}
		_, err = remote_file_storage.S3Write(to, output)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to write to S3: %s", err)
			return
		}
		_, _ = fmt.Fprintf(w, "success: result available at %s", to)
	}

}

func parseParams(w http.ResponseWriter, r *http.Request) (string, string, error) {
	q := r.URL.Query()

	fromValue, queryHasFromParam := q["from"]
	if !queryHasFromParam {
		_, err := fmt.Fprintf(w, "please provide a query param 'from' with a path to the input file")
		return "", "", err
	}
	toValue, queryHasToParam := q["to"]
	if !queryHasToParam {
		_, err := fmt.Fprintf(w, "please provide a query param 'to' with a path to the output file")
		return "", "", err
	}

	return fromValue[0], toValue[0], nil
}

func parseGedcom(inputReader io.Reader, to string) (*[]byte, error) {
	fileScanner := bufio.NewScanner(inputReader)
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

	switch filepath.Ext(to) {
	case ".json":
		return gedcomToJSON(gedcom)
	}

	return nil, fmt.Errorf("no matching file extension in 'to'")
}

func gedcomToJSON(gedcom *model.ConcurrencySafeGedcom) (*[]byte, error) {
	gedcomJson, err := json.Marshal(gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomJson, nil
}
