package expressions

import (
	"banek/ast"
	"banek/exec/operations"
	"strings"
)

type UnaryOp struct {
	Operation operations.UnaryOperator

	Operand ast.Expression
}

func (expr UnaryOp) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expr.Operation.String())
	sb.WriteString(expr.Operand.String())
	sb.WriteByte(')')

	return sb.String()
}

func (expr UnaryOp) IsConst() bool {
	return expr.Operand.IsConst()
}
