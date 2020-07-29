package grpc

import (
	"golang.org/x/net/context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) Parse(_ctx context.Context, paths *PathsToFiles) (*Result, error) {
	log.Printf("parse triggered with paths: %+v\n", paths)
	return &Result{Message: "foo", Error: "bar"}, nil
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
