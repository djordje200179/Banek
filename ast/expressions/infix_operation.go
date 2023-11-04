package expressions

import (
	"banek/ast"
	"banek/exec/operations"
	"strings"
)

type InfixOperation struct {
	Left, Right ast.Expression
	Operator    operations.InfixOperationType
}

func (operation InfixOperation) String() string {
	var sb strings.Builder

	sb.WriteByte('(')
	sb.WriteString(operation.Left.String())
	sb.WriteByte(' ')
	sb.WriteString(operation.Operator.String())
	sb.WriteByte(' ')
	sb.WriteString(operation.Right.String())
	sb.WriteByte(')')

	return sb.String()
}

func (operation InfixOperation) IsConstant() bool {
	return operation.Left.IsConstant() && operation.Right.IsConstant()
}
