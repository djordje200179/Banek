package main

import (
	"banek/compiler"
	"banek/lexer"
	"banek/parser"
	"os"
	"testing"
)

func BenchmarkCompiler(b *testing.B) {
	inputFile, _ := os.Open("test.ba")
	defer inputFile.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tokenChannel := lexer.Tokenize(inputFile, 200)
		statementChannel := parser.Parse(tokenChannel, 20)

		_, _ = compiler.Compile(statementChannel)
	}
}
