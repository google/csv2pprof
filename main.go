package main

import (
	"log"
	"os"

	"github.com/mhansen/csv2pprof"
)

func main() {
	err := csv2pprof.ConvertCSVToCompressedPprof(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
