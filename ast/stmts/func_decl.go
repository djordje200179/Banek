package stmts

import (
	"banek/ast/exprs"
	"banek/symtable"
	"strings"
)

type FuncDecl struct {
	Name      exprs.Ident
	Container *symtable.Container

	Params []exprs.Ident
	Body   Block
}

func (stmt FuncDecl) String() string {
	var sb strings.Builder

	sb.WriteString("func ")
	sb.WriteString(stmt.Name.String())

	sb.WriteByte('(')
	for i, param := range stmt.Params {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(stmt.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (stmt FuncDecl) StmtNode() {}
