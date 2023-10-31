package main

import (
	"banek/bytecode"
	"banek/compiler"
	"banek/lexer"
	"banek/parser"
	"io"
)

func Compile(inputFile io.Reader) (bytecode.Executable, error) {
	tokenChannel := lexer.Tokenize(inputFile, 50)
	statementChannel := parser.Parse(tokenChannel, 5)

	return compiler.Compile(statementChannel)
}
