package main

import (
	"banek/interpreter"
	"banek/interpreter/results"
	"banek/lexer"
	"banek/parser"
	"flag"
	"fmt"
	"os"
)

func HandleUsageError() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: interpreter [file]")
	os.Exit(1)
}

func HandleError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func main() {
	flag.Parse()

	argsCount := flag.NArg()
	if argsCount > 1 {
		HandleUsageError()
	}

	var err error

	var file *os.File
	if argsCount == 1 {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			HandleError(err)
		}

		defer file.Close()
	} else {
		file = os.Stdin
	}

	tokenChannel := lexer.Tokenize(file, 200)
	statementChannel := parser.Parse(tokenChannel, 20)
	resultsChannel := interpreter.Interpret(statementChannel, 5)

	for result := range resultsChannel {
		if err, ok := result.(results.Error); ok {
			HandleError(err)
		}

		fmt.Println(result)
	}
}
