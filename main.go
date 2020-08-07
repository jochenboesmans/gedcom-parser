package main

import (
	"github.com/jochenboesmans/gedcom-parser/grpc"
	"github.com/jochenboesmans/gedcom-parser/parse"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Some functionality might not work if you don't supply environment variables.")
	}
}

func main() {
	checkMainArg()
	switch os.Args[1] {
	case "parse":
		checkFilepathArgs()
		parse.Parse(os.Args[2], os.Args[3])
	case "serve":
		grpc.Serve()
	case "help":
		helpMessage := `
		Usage: 'gedcom-parser <command> <inputFilePath> <outputFilePath>'

		* <command> [REQUIRED]:
			parse - Parse local files. Requires the inputFilePath and outputFilePath to be specified.
			serve - Start a gRPC server for gedcom parsing on remote file storage.

		* <inputFilePath> [OPTIONAL]:
			Relative path to the input file to parse. Please make sure to use the file extensions .ged, .json and .protobuf for respectively GEDCOM, JSON and Protobuf files.
			
		* <outputFilePath> [OPTIONAL]:
			Relative path to the output file. Please make sure to use the file extensions .ged, .json and .protobuf for respectively GEDCOM, JSON and Protobuf files.
		`
		log.Fatal(helpMessage)
	default:
		log.Fatalln("please choose a valid command (use 'gedcom-parser help' for more information on usage)")
	}
}

func checkMainArg() {
	if len(os.Args) < 2 {
		log.Fatalln("please choose a command (use 'gedcom-parser help' for more information on usage)")
	}
}

func checkFilepathArgs() {
	if len(os.Args) < 4 {
		log.Fatalln("please supply inputFilePath and outputFilePath respectively (use 'gedcom-parser help' for more information on usage)")
	}
}
