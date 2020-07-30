package grpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	gedcomSpec "github.com/jochenboesmans/gedcom-parser/gedcom"
	"github.com/jochenboesmans/gedcom-parser/model"
	"github.com/jochenboesmans/gedcom-parser/parse"
	remoteFileStorage "github.com/jochenboesmans/gedcom-parser/remote-file-storage"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) Parse(_ context.Context, paths *PathsToFiles) (*Result, error) {
	log.Printf("started parsing %s to %s", paths.InputFilePath, paths.OutputFilePath)
	log.Printf("reading from s3 bucket at %s...\n", paths.InputFilePath)
	input, err := remoteFileStorage.S3Read(paths.InputFilePath)
	if err != nil {
		errMessage := fmt.Sprintf("failed to read from s3: %s", err)
		log.Println(errMessage)
		return &Result{
			Error: errMessage,
		}, nil
	}

	var output *[]byte
	inputReader := bytes.NewReader(*input)

	switch filepath.Ext(paths.InputFilePath) {
	case ".ged":
		log.Printf("parsing gedcom...\n")
		output, err = parseGedcom(inputReader, paths.OutputFilePath)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse gedcom: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}
	case ".json":
		log.Printf("parsing json...\n")
		output, err = parseJSON(inputReader)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse json: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}
	case ".protobuf":
		log.Printf("parsing protobuf...\n")
		output, err = parseProtobuf(inputReader)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse protobuf: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}

	default:
		errMessage := fmt.Sprintf("failed to match input file extension to: %s", ".ged|.json|.protobuf")
		log.Println(errMessage)
		return &Result{
			Error: errMessage,
		}, nil

	}

	log.Printf("writing to s3 bucket at %s...\n", paths.OutputFilePath)
	_, err = remoteFileStorage.S3Write(paths.OutputFilePath, output)
	if err != nil {
		errMessage := fmt.Sprintf("failed to write to s3: %s", err)
		log.Println(errMessage)
		return &Result{
			Error: errMessage,
		}, nil
	}

	log.Printf("finished parsing %s to %s", paths.InputFilePath, paths.OutputFilePath)
	return &Result{
		Message: fmt.Sprintf("successfully parsed %s to %s", paths.InputFilePath, paths.OutputFilePath),
	}, nil
}

func Serve() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	RegisterParseServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
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
		return parse.GedcomToJSON(gedcom)
	case ".protobuf":
		return parse.GedcomToProto(gedcom)
	}

	return nil, fmt.Errorf("failed to match output file extension to: %s", ".json")
}

func parseJSON(inputReader io.Reader) (*[]byte, error) {
	gedcom := &model.Gedcom{}

	gedcomJson, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	fmt.Println("hello")

	err = json.Unmarshal(gedcomJson, gedcom)
	if err != nil {
		return nil, err
	}

	fmt.Println("hello")

	gedcomBuf := parse.WritableGedcom(gedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}

func parseProtobuf(inputReader io.Reader) (*[]byte, error) {
	var gedcom *model.Gedcom

	gedcomProto, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(gedcomProto, gedcom)
	if err != nil {
		return nil, err
	}

	gedcomBuf := parse.WritableGedcom(gedcom)
	gedcomBytes := gedcomBuf.Bytes()
	return &gedcomBytes, nil
}
