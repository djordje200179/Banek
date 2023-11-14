package main

import (
	"banek/interpreter"
	"banek/interpreter/results"
	"banek/lexer"
	"banek/parser"
	"banek/runtime/types"
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

func REPL() {
	tokenChan := lexer.Tokenize(os.Stdin, 200)
	stmtChan := parser.Parse(tokenChan, 20)
	resultsChan := interpreter.Interpret(stmtChan, 1)

	fmt.Println("Welcome to Banek REPL!")
	fmt.Print(">>> ")

	for result := range resultsChan {
		switch result := result.(type) {
		case results.Error:
			fmt.Println(result)
		case types.Obj:
			fmt.Println(result)
		}

		fmt.Print(">>> ")
	}
}

func EvalFile(file *os.File) {
	tokenChan := lexer.Tokenize(file, 200)
	stmtChan := parser.Parse(tokenChan, 20)
	resultsChan := interpreter.Interpret(stmtChan, 5)

	for result := range resultsChan {
		switch result := result.(type) {
		case results.Error:
			HandleError(result)
		}
	}
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)

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

		EvalFile(file)
	} else {
		REPL()
	}
}
