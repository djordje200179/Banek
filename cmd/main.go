package main

import (
	"banek/lexer"
	"banek/tokens"
	"fmt"
	"os"
)

func main() {
	fileName := "test.ba"
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	l := lexer.New(file)

	for {
		token := l.NextToken()

		if token.Type == tokens.EOF {
			break
		}

		fmt.Println(token)
	}
}
