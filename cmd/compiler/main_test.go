package main

import (
	"banek/analyzer"
	"banek/codegen"
	"banek/emulator"
	"banek/lexer"
	"banek/parser"
	"strings"
	"testing"
)

func runCode(code string) error {
	file := strings.NewReader(code)
	tokens := lexer.Tokenize(file, 100)
	stmts := parser.Parse(tokens, 10)
	stmts = analyzer.Analyze(stmts, 10)
	exec := codegen.Generate(stmts)
	return emulator.Execute(exec)
}

func BenchmarkRecursiveFibonacci(t *testing.B) {
	code := `
		func fibonacci(n) {
		    if n < 2 then
		        return n;
		
		    return fibonacci(n - 1) + fibonacci(n - 2);
		}
	
		println(fibonacci(36));	
	`

	for range t.N {
		err := runCode(code)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkVariableDeclaration(t *testing.B) {
	code := `
		let mut x;
		x = 2;
		println(x);
	`
	for range t.N {
		err := runCode(code)
		if err != nil {
			t.Fatal(err)
		}
	}
}
