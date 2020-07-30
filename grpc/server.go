package grpc

import (
	"bytes"
	"fmt"
	remote_file_storage "github.com/jochenboesmans/gedcom-parser/remote-file-storage"
	"github.com/jochenboesmans/gedcom-parser/server"
	"golang.org/x/net/context"
	"log"
	"net"
	"path/filepath"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) Parse(_ context.Context, paths *PathsToFiles) (*Result, error) {
	log.Printf("started parsing %s to %s", paths.InputFilePath, paths.OutputFilePath)
	log.Printf("reading from s3 bucket at %s...\n", paths.InputFilePath)
	input, err := remote_file_storage.S3Read(paths.InputFilePath)
	if err != nil {
		errMessage := fmt.Sprintf("failed to read from s3: %s", err)
		log.Println(errMessage)
		return &Result{
			Error: errMessage,
		}, nil
	}

	inputReader := bytes.NewReader(*input)
	switch filepath.Ext(paths.InputFilePath) {
	case ".ged":
		log.Printf("parsing gedcom...\n")
		output, err := server.ParseGedcom(inputReader, paths.OutputFilePath)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse gedcom: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}
		log.Printf("writing to s3 bucket at %s...\n", paths.OutputFilePath)
		_, err = remote_file_storage.S3Write(paths.OutputFilePath, output)
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
	errMessage := fmt.Sprintf("failed to match input file extension to: %s", ".ged")
	log.Println(errMessage)
	return &Result{
		Error: errMessage,
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
