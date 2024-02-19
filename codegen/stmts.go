package codegen

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/ast/stmts"
	"banek/bytecode"
	"banek/bytecode/instrs"
)

func (g *generator) compileStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case stmts.Assignment:
		g.compileAssignment(stmt)
	case stmts.CompoundAssignment:
		g.compileCompoundAssignment(stmt)
	case stmts.Block:
		g.compileStmtBlock(stmt)
	case stmts.FuncCall:
		g.compileFuncCallStmt(stmt)
	case stmts.FuncDecl:
		g.compileFuncDecl(stmt)
	case stmts.If:
		g.compileIfStmt(stmt)
	case stmts.Return:
		g.compileReturn(stmt)
	case stmts.VarDecl:
		g.compileVarDecl(stmt)
	case stmts.While:
		g.compileWhile(stmt)
	default:
		panic("unreachable")
	}
}

func (g *generator) compileAssignment(stmt stmts.Assignment) {

}

func (g *generator) compileCompoundAssignment(stmt stmts.CompoundAssignment) {

}

func (g *generator) compileStmtBlock(stmt stmts.Block) {
	for _, stmt := range stmt {
		g.compileStmt(stmt)
	}
}

func (g *generator) compileFuncDecl(stmt stmts.FuncDecl) {
	if g.level == 0 {
		g.vars++
	}

	funcIndex := len(g.funcPool)
	g.container = &container{
		level:    g.container.level + 1,
		index:    funcIndex,
		previous: g.container,
		vars:     len(stmt.Params),
	}

	g.compileStmtBlock(stmt.Body)
	g.emitInstr(instrs.OpPushUndef)
	g.emitInstr(instrs.OpReturn)

	g.funcPool = append(g.funcPool, bytecode.FuncTemplate{
		Name:      stmt.Name.String(),
		NumParams: len(stmt.Params),
		NumLocals: g.container.vars,
		Code:      g.container.code,
	})

	g.container = g.container.previous

	g.emitInstr(instrs.OpMakeFunc, funcIndex)
	g.compileStore(stmt.Name)
}

func (g *generator) compileFuncCallStmt(expr stmts.FuncCall) {
	g.compileFuncCall(exprs.FuncCall(expr))
	g.emitInstr(instrs.OpPop)
}

func (g *generator) compileIfStmt(stmt stmts.If) {
	g.compileExpr(stmt.Cond)
	jumpPc := g.container.currAddr()
	g.emitInstr(instrs.OpBranchFalse, 0)

	g.compileStmt(stmt.Cons)

	if stmt.Alt != nil {
		altJumpPc := g.container.currAddr()
		g.emitInstr(instrs.OpJump, 0)
		g.container.patchJumpOperand(jumpPc, 0)
		g.compileStmt(stmt.Alt)
		g.container.patchJumpOperand(altJumpPc, 0)
	} else {
		g.container.patchJumpOperand(jumpPc, 0)
	}
}

func (g *generator) compileReturn(stmt stmts.Return) {
	if stmt.Value != nil {
		g.compileExpr(stmt.Value)
	}

	g.emitInstr(instrs.OpReturn)
}

func (g *generator) compileVarDecl(stmt stmts.VarDecl) {
	if g.level == 0 {
		g.vars++
	}
}

func (g *generator) compileWhile(stmt stmts.While) {

}
