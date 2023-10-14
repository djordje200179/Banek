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
	l := lexer.New(in)

	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(program.String())
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
