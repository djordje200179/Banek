package compiler

import (
	"banek/ast"
	"banek/bytecode"
	"banek/compiler/scopes"
	"banek/exec/objects"
	"slices"
)

type compiler struct {
	constants []objects.Object
	functions []bytecode.FunctionTemplate

	globalScope scopes.Global
	scopes      []scopes.Scope
}

func Compile(statementsChannel <-chan ast.Statement) (bytecode.Executable, error) {
	compiler := &compiler{
		scopes: make([]scopes.Scope, 1),
	}
	compiler.scopes[0] = &compiler.globalScope

	for statement := range statementsChannel {
		err := compiler.compileStatement(statement)
		if err != nil {
			return bytecode.Executable{}, err
		}
	}

	return compiler.makeExecutable(), nil
}

func (compiler *compiler) addConstant(object objects.Object) int {
	if index := slices.Index(compiler.constants, object); index != -1 {
		return index
	}

	index := len(compiler.constants)
	compiler.constants = append(compiler.constants, object)

	return index
}

func (compiler *compiler) addFunction(template bytecode.FunctionTemplate) int {
	index := len(compiler.functions)
	compiler.functions = append(compiler.functions, template)

	return index
}

func (compiler *compiler) topScope() scopes.Scope {
	return compiler.scopes[len(compiler.scopes)-1]
}

func (compiler *compiler) popScope() {
	compiler.scopes = compiler.scopes[:len(compiler.scopes)-1]
}

func (compiler *compiler) pushScope(scope scopes.Scope) {
	compiler.scopes = append(compiler.scopes, scope)
}

func (compiler *compiler) makeExecutable() bytecode.Executable {
	executable := compiler.globalScope.MakeExecutable()
	executable.ConstantsPool = compiler.constants
	executable.FunctionsPool = compiler.functions

	return executable
}
