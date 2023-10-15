package statements

import (
	"banek/ast"
	"strings"
)

type If struct {
	Condition ast.Expression

	Consequence, Alternative ast.Statement
}

func (statement If) StatementNode() {}

func (statement If) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(statement.Condition.String())
	sb.WriteString(" then {\n")
	sb.WriteString(statement.Consequence.String())
	sb.WriteString("\n}")
	if statement.Alternative != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(statement.Alternative.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}
