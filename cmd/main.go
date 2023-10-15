package main

import (
	"banek/interpreter"
	"banek/interpreter/objects"
	"banek/interpreter/results"
	"banek/lexer"
	"banek/parser"
	"fmt"
	"os"
)

func main() {
	fileName := "test.ba"
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tokenChannel := lexer.New(file).Tokenize(100)
	statementChannel := parser.New().Parse(tokenChannel, 100)
	objectsChannel := interpreter.New().Eval(statementChannel, 100)

	for object := range objectsChannel {
		switch object := object.(type) {
		case results.Error:
			_, _ = fmt.Fprintln(os.Stderr, object.Error())
		case objects.Object:
			_, _ = fmt.Fprintln(os.Stdout, object.String())
		}
	}
}
