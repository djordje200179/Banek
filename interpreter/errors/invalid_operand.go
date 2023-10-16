package errors

import (
	"banek/interpreter/objects"
	"fmt"
)

type InvalidOperandError struct {
	// TODO: Better error structure and message
	Operator string
	Operand  objects.Object
}

func (err InvalidOperandError) Error() string {
	return fmt.Sprintf("invalid operand: expected %s, got %s", err.Operator, err.Operand.Type())
}
