package main

import (
	"banek/analyzer"
	"banek/codegen"
	vm "banek/emulator"
	"banek/lexer"
	"banek/parser"
	"fmt"
	"strings"
)

const sampleProgram = `
	func fibonacci(n) {
	    if n < 2 then
	        return n;
	
	    return fibonacci(n - 1) + fibonacci(n - 2);
	}

	println(fibonacci(20));	
`

func main() {
	file := strings.NewReader(sampleProgram)
	tokens := lexer.Tokenize(file, 100)
	stmts := parser.Parse(tokens, 10)
	stmts = analyzer.Analyze(stmts, 10)

	exec := codegen.Generate(stmts)

	fmt.Println(exec.String())

	vm.Execute(exec)
}
