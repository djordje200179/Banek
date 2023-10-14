package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type PrefixOperation struct {
	Operator tokens.Token
	Operand  ast.Expression
}

func (expression PrefixOperation) ExpressionNode() {}

func (expression PrefixOperation) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expression.Operator.String())
	sb.WriteString(expression.Operand.String())
	sb.WriteByte(')')

	return sb.String()
}
