package main

import (
	"banek/interpreter"
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

	tokenChannel := lexer.Tokenize(file, 50)
	statementChannel := parser.Parse(tokenChannel, 5)
	resultsChannel := interpreter.Interpret(statementChannel, 5)

	for result := range resultsChannel {
		fmt.Println(result)
	}
}
