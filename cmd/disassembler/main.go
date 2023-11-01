package main

import (
	"banek/bytecode"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func HandleUsageError() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: disassembler <file>")
	flag.PrintDefaults()

	os.Exit(1)
}

func HandleError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)

	flag.Parse()

	argsCount := flag.NArg()
	if argsCount != 1 {
		HandleUsageError()
	}

	inputFile, err := os.Open(flag.Arg(0))
	if err != nil {
		HandleError(err)
	}

	defer inputFile.Close()

	var executable bytecode.Executable
	err = gob.NewDecoder(inputFile).Decode(&executable)
	if err != nil {
		HandleError(err)
	}

	fmt.Println(executable)
}
