package compiler

import (
	"banek/ast"
	"banek/bytecode"
	"banek/compiler/scopes"
	"banek/runtime/types"
	"slices"
)

type compiler struct {
	consts []types.Obj
	funcs  []bytecode.FuncTemplate

	globalScope scopes.Global
	scopes      []scopes.Scope
}

func Compile(stmtsChan <-chan ast.Stmt) (bytecode.Executable, error) {
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

func (compiler *compiler) addConst(object types.Obj) int {
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
	scopeStackLen := len(compiler.scopes)

	if blockScope, ok := compiler.scopes[scopeStackLen-1].(*scopes.Block); ok {
		compiler.scopes[scopeStackLen-1] = blockScope.Parent
		return
	}

	compiler.scopes = compiler.scopes[:scopeStackLen-1]
}

func (compiler *compiler) pushScope(scope scopes.Scope) {
	if blockScope, ok := scope.(*scopes.Block); ok {
		compiler.scopes[len(compiler.scopes)-1] = blockScope
		return
	}

	compiler.scopes = append(compiler.scopes, scope)
}

func (compiler *compiler) makeExecutable() bytecode.Executable {
	executable := compiler.globalScope.MakeExecutable()
	executable.ConstsPool = compiler.consts
	executable.FuncsPool = compiler.funcs

	return executable
}
