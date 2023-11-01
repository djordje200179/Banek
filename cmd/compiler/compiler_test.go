package main

import (
	"banek/compiler"
	"banek/lexer"
	"banek/parser"
	"os"
	"testing"
)

func BenchmarkCompiler(b *testing.B) {
	inputFile, err := os.Open("test.ba")
	if err != nil {
		b.Fatal(err)
	}
	defer inputFile.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		inputFile.Seek(0, 0)

		tokenChannel := lexer.Tokenize(inputFile, 200)
		statementChannel := parser.Parse(tokenChannel, 20)

		_, err := compiler.Compile(statementChannel)
		if err != nil {
			b.Error(err)
		}
	}
}
