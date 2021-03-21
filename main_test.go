package main

import (
	"github.com/jochenboesmans/gedcom-parser/parse"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func BenchmarkParseITISGedToJson(b *testing.B) {
	cpuFile, err := os.Create("cpu-itis.prof")
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		return
	}
	parse.Parse("examples/ITIS.ged", "test-output/ITIS.json")
	pprof.StopCPUProfile()

	memFile, err := os.Create("mem-itis.prof")
	if err != nil {
		return
	}
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		return
	}

}

func BenchmarkParseHPGedToJson(b *testing.B) {
	cpuFile, err := os.Create("cpu-hp.prof")
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		return
	}
	parse.Parse("examples/harry_potter.ged", "test-output/harry_potter.json")
	pprof.StopCPUProfile()

	memFile, err := os.Create("mem-hp.prof")
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		return
	}

}

func BenchmarkParseWikiGodsGedToJson(b *testing.B) {
	cpuFile, err := os.Create("cpu-wg.prof")
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		return
	}
	parse.Parse("examples/wikipedia_gods.ged", "test-output/wikipedia_gods.json")
	pprof.StopCPUProfile()

	memFile, err := os.Create("mem-wg.prof")
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		return
	}

}
