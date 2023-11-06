package statements

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type VarDecl struct {
	Name  expressions.Identifier
	Value ast.Expression

	Mutable bool
}

func (stmt VarDecl) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Let.String())
	sb.WriteByte(' ')

	if stmt.Mutable {
		sb.WriteString(tokens.Mut.String())
		sb.WriteByte(' ')
	}

	sb.WriteString(stmt.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(stmt.Value.String())

	return sb.String()
}

func (stmt VarDecl) HasSideEffects() bool {
	return !stmt.Value.IsConst()
}
