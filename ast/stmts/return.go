package stmts

import (
	"banek/ast"
	"strings"
)

type Return struct {
	Value ast.Expr
}

func (stmt Return) String() string {
	var sb strings.Builder

	sb.WriteString("return ")
	sb.WriteString(stmt.Value.String())

	return sb.String()
}

func (stmt Return) StmtNode() {}
