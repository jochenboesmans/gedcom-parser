package main

import (
	"github.com/jochenboesmans/gedcom-parser/grpc"
	"github.com/jochenboesmans/gedcom-parser/parse"
	"os"
)

//var inputFilePath = flag.String("input", "", "path to input file")
//var outputFilePath = flag.String("output", "", "path to output file")

func main() {
	switch os.Args[1] {
	case "parse":
		parse.Parse(os.Args[2], os.Args[3])
	case "serve":
		grpc.Serve()
	}
}
