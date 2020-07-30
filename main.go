package main

import (
	"github.com/jochenboesmans/gedcom-parser/grpc"
	"github.com/jochenboesmans/gedcom-parser/parse"
	"os"
)

func main() {
	switch os.Args[1] {
	case "parse":
		parse.Parse()
	case "serve":
		grpc.Serve()
	}
}
