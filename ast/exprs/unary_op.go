package exprs

import (
	"banek/ast"
	"banek/runtime/ops"
	"strings"
)

type UnaryOp struct {
	Operator ops.UnaryOperator

	Operand ast.Expr
}

func (expr UnaryOp) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expr.Operator.String())
	sb.WriteString(expr.Operand.String())
	sb.WriteByte(')')

	return sb.String()
}

func (expr UnaryOp) IsConst() bool {
	return expr.Operand.IsConst()
}
