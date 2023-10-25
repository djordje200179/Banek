package errors

import (
	"banek/exec/objects"
	"banek/tokens"
	"strings"
)

type ErrUnknownOperator struct {
	Operator tokens.TokenType
}

func (err ErrUnknownOperator) Error() string {
	return "unknown operator: " + err.Operator.String()
}

type ErrInvalidOperand struct {
	Operation string

	LeftOperand  objects.Object
	RightOperand objects.Object
}

func (err ErrInvalidOperand) Error() string {
	var sb strings.Builder

	switch {
	case err.LeftOperand != nil && err.RightOperand != nil:
		sb.WriteString("invalid operands for ")
		sb.WriteString(err.Operation)
		sb.WriteString(": ")
		sb.WriteString(err.LeftOperand.Type())
		sb.WriteString(" and ")
		sb.WriteString(err.RightOperand.Type())
	case err.LeftOperand != nil:
		sb.WriteString("invalid operand for ")
		sb.WriteString(err.Operation)
		sb.WriteString(": ")
		sb.WriteString(err.LeftOperand.Type())
	case err.RightOperand != nil:
		sb.WriteString("invalid operand for ")
		sb.WriteString(err.Operation)
		sb.WriteString(": ")
		sb.WriteString(err.RightOperand.Type())
	}

	return sb.String()
}
