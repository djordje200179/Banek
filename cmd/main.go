package main

import (
	"banek/compiler"
	"banek/lexer"
	"banek/parser"
	"banek/vm"
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

	executable, err := compiler.Compile(statementChannel)
	if err != nil {
		panic(err)
	}

	err = vm.Execute(executable)
	if err != nil {
		panic(err)
	}
}
