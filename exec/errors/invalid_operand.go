package errors

import (
	"banek/exec/objects"
	"strings"
)

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
