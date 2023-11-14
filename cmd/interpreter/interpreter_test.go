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
	inputFile, err := os.Open("eratosthenes_sieve.ba")
	if err != nil {
		b.Fatal(err)
	}
	defer inputFile.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		inputFile.Seek(0, 0)

		tokenChan := lexer.Tokenize(inputFile, 200)
		stmtChan := parser.Parse(tokenChan, 20)
		resultsChan := interpreter.Interpret(stmtChan, 5)

		for result := range resultsChan {
			if err, ok := result.(results.Error); ok {
				b.Error(err)
			}
		}
	}
}
