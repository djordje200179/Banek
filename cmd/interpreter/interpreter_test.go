package main

import (
	"banek/interpreter"
	"banek/interpreter/results"
	"banek/lexer"
	"banek/parser"
	"os"
	"testing"
)

func BenchmarkInterpreter(b *testing.B) {
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
		resultsChannel := interpreter.Interpret(statementChannel, 5)

		for result := range resultsChannel {
			if err, ok := result.(results.Error); ok {
				b.Error(err)
			}
		}
	}
}
