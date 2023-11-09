package builtins

import (
	"banek/runtime/objs"
	"banek/runtime/types"
)

func builtinLen(args []types.Obj) (types.Obj, error) {
	if len(args) != 1 {
		return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
	}

	switch arg := args[0].(type) {
	case objs.Str:
		return objs.Int(len(arg)), nil
	case objs.Array:
		return objs.Int(len(arg)), nil
	default:
		return objs.Undefined{}, nil
	}
}
