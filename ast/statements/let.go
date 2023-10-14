package statements

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type LetStatement struct {
	Name  expressions.Identifier
	Value ast.Expression
}

func (statement LetStatement) StatementNode() {}

func (statement LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Let.String())
	sb.WriteByte(' ')
	sb.WriteString(statement.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(statement.Value.String())

	return sb.String()
}
