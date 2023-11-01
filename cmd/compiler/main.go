package main

import (
	"banek/compiler"
	"banek/lexer"
	"banek/parser"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func HandleUsageError() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: compiler <file>")
	flag.PrintDefaults()

	os.Exit(1)
}

func HandleError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)

	var outputFileName string
	flag.StringVar(&outputFileName, "o", "a.bac", "Output file name")

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

	tokenChannel := lexer.Tokenize(inputFile, 200)
	statementChannel := parser.Parse(tokenChannel, 20)

	executable, err := compiler.Compile(statementChannel)
	if err != nil {
		HandleError(err)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		HandleError(err)
	}

	defer outputFile.Close()

	err = gob.NewEncoder(outputFile).Encode(executable)
	if err != nil {
		HandleError(err)
	}
}
