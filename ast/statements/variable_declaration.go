package statements

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type VariableDeclaration struct {
	Name  expressions.Identifier
	Value ast.Expression

	Mutable bool
}

func (statement VariableDeclaration) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Let.String())
	sb.WriteByte(' ')

	if statement.Mutable {
		sb.WriteString(tokens.Mut.String())
		sb.WriteByte(' ')
	}

	sb.WriteString(statement.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(statement.Value.String())

	return sb.String()
}

func (statement VariableDeclaration) HasSideEffects() bool {
	return !statement.Value.IsConstant()
}
