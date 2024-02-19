package main

import (
	"banek/analyzer"
	"banek/codegen"
	"banek/emulator"
	"banek/lexer"
	"banek/parser"
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

	tokens := lexer.Tokenize(inputFile, 100)
	stmts := parser.Parse(tokens, 10)
	stmts = analyzer.Analyze(stmts, 10)
	exec := codegen.Generate(stmts)

	err = emulator.Execute(exec)
	if err != nil {
		HandleError(err)
	}
}
