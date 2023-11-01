package main

import (
	"banek/bytecode"
	"banek/vm"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func HandleUsageError() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: emulator <file>")
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

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		HandleError(err)
	}
	defer file.Close()

	var executable bytecode.Executable
	err = gob.NewDecoder(file).Decode(&executable)
	if err != nil {
		HandleError(err)
	}

	err = vm.Execute(executable)
	if err != nil {
		HandleError(err)
	}
}
