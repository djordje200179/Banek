package builtins

import (
	"banek/runtime/objs"
	"fmt"
)

type InvalidTypeError struct {
	ArgIndex int

	Arg objs.Obj
}

func (err InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type for %d. argument: %s", err.ArgIndex+1, err.Arg.String())
}

type CallError struct {
	BuiltinName string
}

func (err CallError) Error() string {
	return fmt.Sprintf("call to builtin %s failed", err.BuiltinName)
}
