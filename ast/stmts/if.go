package stmts

import (
	"banek/ast"
	"strings"
)

type If struct {
	Cond ast.Expr

	Cons, Alt ast.Stmt
}

func (stmt If) String() string {
	var sb strings.Builder

	sb.WriteString("if ")
	sb.WriteString(stmt.Cond.String())
	sb.WriteString(" then {\n")
	sb.WriteString(stmt.Cons.String())
	sb.WriteString("\n}")
	if stmt.Alt != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(stmt.Alt.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}

func (stmt If) StmtNode() {}
