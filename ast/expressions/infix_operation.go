package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type InfixOperation struct {
	Left, Right ast.Expression
	Operator    tokens.Token
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
