package grpc

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jochenboesmans/gedcom-parser/parse"
	remoteFileStorage "github.com/jochenboesmans/gedcom-parser/remote-file-storage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	ParseServiceServer
}

func (s *Server) init() {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))))

	s.uploader = s3manager.NewUploader(sess)
	s.downloader = s3manager.NewDownloader(sess)
}

func (s *Server) Parse(_ context.Context, paths *PathsToFiles) (*Result, error) {
	log.Printf("started parsing %s to %s", paths.InputFilePath, paths.OutputFilePath)
	log.Printf("reading from s3 bucket at %s...\n", paths.InputFilePath)
	input, err := remoteFileStorage.S3Read(paths.InputFilePath, s.downloader)
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
		output, err = parse.ParseGedcom(inputReader, paths.OutputFilePath)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse gedcom: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}
	case ".json":
		log.Printf("parsing json...\n")
		output, err = parse.ParseJSON(inputReader)
		if err != nil {
			errMessage := fmt.Sprintf("failed to parse json: %s", err)
			log.Println(errMessage)
			return &Result{
				Error: errMessage,
			}, nil
		}
	case ".protobuf":
		log.Printf("parsing protobuf...\n")
		output, err = parse.ParseProtobuf(inputReader)
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
	_, err = remoteFileStorage.S3Write(paths.OutputFilePath, output, s.uploader)
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
	s.init()

	grpcServer := grpc.NewServer()

	RegisterParseServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
