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
	
		println(fibonacci(40));	
	`

	for range t.N {
		err := runCode(code)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestVariableDeclaration(t *testing.T) {
	code := `
		let mut x;
		x = 2;
		println(x);
	`
	err := runCode(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayManipulation(t *testing.T) {
	code := `
		let arr = [1, 2, 3, 4, 5];
		arr[1] += 1;
		println(arr[1]);
		println([0] + arr);
	`
	err := runCode(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAnonymousFunctions(t *testing.T) {
	code := `
		let f = |x| -> x + 1;
		println(f(5));
	`
	err := runCode(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWhileLoop(t *testing.T) {
	code := `
		let mut x = 0;
		while x < 5 do {
			x += 1;
			print(x);
		}
		println();
		println("Final:", x);
	`
	err := runCode(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestForLoop(t *testing.T) {
	code := `
		let mut i;
		for i = 0; i < 5; i += 1 do
			print(i);
	`
	err := runCode(code)
	if err != nil {
		t.Fatal(err)
	}
}
