package main

import (
	"banek/interpreter"
	"banek/lexer"
	"banek/parser"
	"os"
	"testing"
)

func BenchmarkInterpreter(b *testing.B) {
	inputFile, _ := os.Open("test.ba")
	defer inputFile.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tokenChannel := lexer.Tokenize(inputFile, 200)
		statementChannel := parser.Parse(tokenChannel, 20)
		resultsChannel := interpreter.Interpret(statementChannel, 5)

		for range resultsChannel {
			<-resultsChannel
		}
	}
}
