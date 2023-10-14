package statements

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type VariableDeclarationStatement struct {
	Name  expressions.Identifier
	Value ast.Expression

	Const bool
}

func (statement VariableDeclarationStatement) StatementNode() {}

func (statement VariableDeclarationStatement) String() string {
	var sb strings.Builder

	if statement.Const {
		sb.WriteString(tokens.Const.String())
	} else {
		sb.WriteString(tokens.Let.String())
	}
	sb.WriteByte(' ')
	sb.WriteString(statement.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(statement.Value.String())

	return sb.String()
}
