package stmts

import (
	"banek/ast"
	"strings"
)

type While struct {
	Cond ast.Expr

	Body ast.Stmt
}

func (stmt While) String() string {
	var sb strings.Builder

	sb.WriteString("while")
	sb.WriteString(stmt.Cond.String())
	sb.WriteString(" do {\n")
	sb.WriteString(stmt.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (stmt While) StmtNode() {}
