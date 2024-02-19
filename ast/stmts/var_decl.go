package stmts

import (
	"banek/ast"
	"banek/ast/exprs"
	"strings"
)

type VarDecl struct {
	Var   exprs.Ident
	Value ast.Expr

	Mutable bool
}

func (stmt VarDecl) String() string {
	var sb strings.Builder

	sb.WriteString("let ")
	if stmt.Mutable {
		sb.WriteString("mut ")
	}

	sb.WriteString(stmt.Var.String())
	sb.WriteString(" = ")
	sb.WriteString(stmt.Value.String())

	return sb.String()
}

func (stmt VarDecl) StmtNode() {}
