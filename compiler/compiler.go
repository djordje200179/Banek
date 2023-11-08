package compiler

import (
	"banek/ast"
	"banek/bytecode"
	"banek/compiler/scopes"
	"banek/exec/objects"
	"slices"
)

type compiler struct {
	consts []objects.Object
	funcs  []bytecode.FuncTemplate

	globalScope scopes.Global
	scopes      []scopes.Scope
}

func Compile(stmtsChan <-chan ast.Statement) (bytecode.Executable, error) {
	compiler := &compiler{
		scopes: make([]scopes.Scope, 1),
	}
	compiler.scopes[0] = &compiler.globalScope

	for stmt := range stmtsChan {
		err := compiler.compileStmt(stmt)
		if err != nil {
			return bytecode.Executable{}, err
		}
	}

	return compiler.makeExecutable(), nil
}

func (compiler *compiler) addConst(object objects.Object) int {
	if index := slices.IndexFunc(compiler.consts, object.Equals); index != -1 {
		return index
	}

	index := len(compiler.consts)
	compiler.consts = append(compiler.consts, object)

	return index
}

func (compiler *compiler) addFunc(template bytecode.FuncTemplate) int {
	index := len(compiler.funcs)
	compiler.funcs = append(compiler.funcs, template)

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
	executable.ConstsPool = compiler.consts
	executable.FuncsPool = compiler.funcs

	return executable
}
