// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	stackSep = flag.String("stacksep", ";", "Character to split the 'stack' column into lines with")
)

func openInput() io.Reader {
	args := flag.Args()
	if len(args) == 0 {
		return os.Stdin
	}
	in, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	return in
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "csv2pprof < input.csv > output.pprof.gz\n")
		fmt.Fprint(os.Stderr, "csv2pprof input.csv > output.pprof.gz\n")
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Parameters:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	in := openInput()
	err := ConvertCSVToCompressedPprof(in, os.Stdout, *stackSep)
	if err != nil {
		log.Fatal(err)
	}
}
