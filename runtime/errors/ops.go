package errors

import (
	"banek/runtime/types"
	"strings"
)

type ErrUnknownOperator struct {
	Operator string
}

func (err ErrUnknownOperator) Error() string {
	return "unknown operator: " + err.Operator
}

type ErrInvalidOp struct {
	Operator string

	LeftOperand  types.Obj
	RightOperand types.Obj
}

func (err ErrInvalidOp) Error() string {
	var sb strings.Builder

	switch {
	case err.LeftOperand != nil && err.RightOperand != nil:
		sb.WriteString("invalid operands for ")
		sb.WriteString(err.Operator)
		sb.WriteString(": ")
		sb.WriteString(err.LeftOperand.String())
		sb.WriteString(" and ")
		sb.WriteString(err.RightOperand.String())
	case err.LeftOperand != nil:
		sb.WriteString("invalid operand for ")
		sb.WriteString(err.Operator)
		sb.WriteString(": ")
		sb.WriteString(err.LeftOperand.String())
	case err.RightOperand != nil:
		sb.WriteString("invalid operand for ")
		sb.WriteString(err.Operator)
		sb.WriteString(": ")
		sb.WriteString(err.RightOperand.String())
	}

	return sb.String()
}
