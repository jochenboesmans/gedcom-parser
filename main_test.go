package main

import (
	"github.com/jochenboesmans/gedcom-parser/parse"
	"os"
	"runtime/pprof"
	"testing"
)

func BenchmarkParseITISGedToJson(b *testing.B) {
	file, err := os.Create("./heap-profile.pb.gz")
	if err != nil {
		return
	}
	err = pprof.WriteHeapProfile(file)
	if err != nil {
		return
	}
	b.ReportAllocs()
	parse.Parse("examples/ITIS.ged", "test-output/ITIS.json")
	err = pprof.WriteHeapProfile(file)
	if err != nil {
		return
	}
}

func BenchmarkParseITISJsonToGed(b *testing.B) {
	//b.ReportAllocs()
	//parse.Parse("test-output/ITIS.json", "test-output/ITIS.ged")
}
