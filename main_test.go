package main

import (
	"github.com/jochenboesmans/gedcom-parser/parse"
	"os"
	"runtime/pprof"
	"testing"
)

func BenchmarkParseITISGedToJson(b *testing.B) {
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		return
	}
	parse.Parse("examples/ITIS.ged", "test-output/ITIS.json")
	pprof.StopCPUProfile()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		return
	}
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		return
	}

}

func BenchmarkParseITISJsonToGed(b *testing.B) {
	//b.ReportAllocs()
	//parse.Parse("test-output/ITIS.json", "test-output/ITIS.ged")
}
