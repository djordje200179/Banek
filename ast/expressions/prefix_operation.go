package expressions

import (
	"banek/ast"
	"banek/exec/operations"
	"strings"
)

type PrefixOperation struct {
	Operator operations.PrefixOperationType
	Operand  ast.Expression
}

func (expression PrefixOperation) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(expression.Operator.String())
	sb.WriteString(expression.Operand.String())
	sb.WriteByte(')')

	return sb.String()
}

func (expression PrefixOperation) IsConstant() bool {
	return expression.Operand.IsConstant()
}
