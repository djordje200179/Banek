package main

import (
	"banek/evaluator"
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
	objectsChannel := evaluator.New().Evaluate(statementChannel, 100)

	for object := range objectsChannel {
		fmt.Println(object.String())
	}
}
