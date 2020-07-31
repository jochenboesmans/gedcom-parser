package main

import (
	"github.com/jochenboesmans/gedcom-parser/grpc"
	"github.com/jochenboesmans/gedcom-parser/parse"
	"log"
	"os"
)

//var inputFilePath = flag.String("input", "", "path to input file")
//var outputFilePath = flag.String("output", "", "path to output file")

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please choose a command (parse|serve)")
	}
	switch os.Args[1] {
	case "parse":
		if len(os.Args) < 4 {
			log.Fatalln("please supply inputFilePath and outputFilePath respectively")
		}
		parse.Parse(os.Args[2], os.Args[3])
	case "serve":
		grpc.Serve()
	}
}
