package main

import (
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile:", err)
	}
	defer f.Close()

	fh, err := os.Create("heap.prof")
	if err != nil {
		log.Fatal("could not create CPU profile:", err)
	}
	defer fh.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile:", err)
	}
	defer pprof.StopCPUProfile()

	if err := pprof.WriteHeapProfile(fh); err != nil {
		log.Fatal("could not write heap profile:", err)
	}

	runCLI(os.Args)
}
