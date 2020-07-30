package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestServe(t *testing.T) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewParseServiceClient(conn)

	response, err := c.Parse(context.Background(), &PathsToFiles{InputFilePath: "json/ITIS.json", OutputFilePath: "gedcom/new-ITIS.ged"})
	if err != nil {
		log.Fatalf("error when calling Parse: %s", err)
	}
	log.Printf("Response from server: %+v", response)
}
