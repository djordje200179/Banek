package builtins

import (
	"banek/runtime"
	"banek/runtime/primitives"
)

func builtinLen(args []runtime.Obj) (runtime.Obj, error) {
	switch arg := args[0].(type) {
	case primitives.String:
		return primitives.Int(len(arg)), nil
	case primitives.Array:
		return primitives.Int(len(arg)), nil
	default:
		// TODO: error
		return nil, nil
	}
}
