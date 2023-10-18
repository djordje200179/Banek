package errors

import (
	"banek/exec/objects"
	"fmt"
)

type ErrInvalidOperand struct {
	// TODO: Better error structure and message
	Operator string
	Operand  objects.Object
}

func (err ErrInvalidOperand) Error() string {
	return fmt.Sprintf("invalid operand: expected %s, got %s", err.Operator, err.Operand.Type())
}
