package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type PrefixedExpression struct {
	Operator tokens.Token
	Wrapped  ast.Expression
}

func (expression PrefixedExpression) ExpressionNode() {}

func (expression PrefixedExpression) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expression.Operator.String())
	sb.WriteString(expression.Wrapped.String())
	sb.WriteByte(')')

	return sb.String()
}
