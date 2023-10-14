package main

import (
	"banek/lexer"
	"banek/parser"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	tokenChannel := lexer.New(in).Tokenize(100)
	statementChannel := parser.New().Parse(tokenChannel, 100)

	for statement := range statementChannel {
		if statement.Error != nil {
			fmt.Println(statement.Error)
			continue
		}

		fmt.Println(statement.Statement.String())
	}
}

func main() {
	fileName := "test.ba"
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	Start(file, os.Stdout)
}
