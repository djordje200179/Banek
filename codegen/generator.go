package codegen

import (
	"banek/ast"
	"banek/bytecode"
	"banek/bytecode/instrs"
)

func Generate(stmtChan <-chan ast.Stmt) *bytecode.Executable {
	g := &generator{
		funcPool: make([]bytecode.FuncTemplate, 1),
	}
	g.container = &g.global

	for stmt := range stmtChan {
		g.compileStmt(stmt)
	}

	g.emitInstr(instrs.OpHalt)

	g.funcPool[0] = bytecode.FuncTemplate{
		Name:       "<entry>",
		NumLocals:  g.global.vars,
		Code:       g.global.code,
		IsCaptured: true,
	}

	return g.makeExecutable()
}

type generator struct {
	global container
	*container

	stringPool []string
	funcPool   []bytecode.FuncTemplate
}

func (g *generator) makeExecutable() *bytecode.Executable {
	return &bytecode.Executable{
		StringPool: g.stringPool,
		FuncPool:   g.funcPool,
	}
}
