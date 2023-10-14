package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type InfixExpression struct {
	Left, Right ast.Expression
	Operator    tokens.Token
}

func (expression InfixExpression) ExpressionNode() {}

func (expression InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expression.Left.String())
	sb.WriteByte(' ')
	sb.WriteString(expression.Operator.String())
	sb.WriteByte(' ')
	sb.WriteString(expression.Right.String())
	sb.WriteByte(')')

	return sb.String()
}
