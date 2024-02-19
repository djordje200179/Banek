package exprs

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type BinaryOp struct {
	Left, Right ast.Expr

	Operator tokens.Type
}

func (expr BinaryOp) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expr.Left.String())
	sb.WriteByte(' ')
	sb.WriteString(expr.Operator.String())
	sb.WriteByte(' ')
	sb.WriteString(expr.Right.String())
	sb.WriteByte(')')

	return sb.String()
}

func (expr BinaryOp) IsConst() bool {
	return expr.Left.IsConst() && expr.Right.IsConst()
}
