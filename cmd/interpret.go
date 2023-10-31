package main

import (
	"banek/interpreter"
	"banek/interpreter/results"
	"banek/lexer"
	"banek/parser"
	"fmt"
	"io"
)

func Interpret(inputFile io.Reader) error {
	tokenChannel := lexer.Tokenize(inputFile, 500)
	statementChannel := parser.Parse(tokenChannel, 50)
	resultsChannel := interpreter.Interpret(statementChannel, 5)

	for result := range resultsChannel {
		if err, ok := result.(results.Error); ok {
			return err
		}

		fmt.Println(result)
	}

	return nil
}
