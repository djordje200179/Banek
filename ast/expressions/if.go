package expressions

import (
	"banek/ast"
	"strings"
)

type If struct {
	Condition ast.Expression

	Consequence, Alternative ast.Expression
}

func (expression If) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(expression.Condition.String())
	sb.WriteString(" then {\n")
	sb.WriteString(expression.Consequence.String())
	sb.WriteString("\n}")
	if expression.Alternative != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(expression.Alternative.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}

func (expression If) IsConstant() bool {
	return expression.Condition.IsConstant() && expression.Consequence.IsConstant() && expression.Alternative.IsConstant()
}
