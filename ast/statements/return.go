package statements

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type ReturnStatement struct {
	Value ast.Expression
}

func (statement ReturnStatement) StatementNode() {}

func (statement ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Return.String())
	sb.WriteByte(' ')
	sb.WriteString(statement.Value.String())

	return sb.String()
}
