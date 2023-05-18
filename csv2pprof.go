package main

import (
	"log"
	"os"
)

func main() {
	err := ConvertCSVToCompressedPprof(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
