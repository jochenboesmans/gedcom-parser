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
		log.Println("No .env file found. Some functionality might not work.")
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
	default:
		log.Fatalln("please choose a command in (parse|serve)")
	}
}

func checkMainArg() {
	if len(os.Args) < 2 {
		log.Fatalln("please choose a command (parse|serve)")
	}
}

func checkFilepathArgs() {
	if len(os.Args) < 4 {
		log.Fatalln("please supply inputFilePath and outputFilePath respectively")
	}
}
