package builtins

import (
	"banek/runtime/objs"
	"banek/runtime/types"
	"strconv"
)

func builtinStr(args []types.Obj) (types.Obj, error) {
	if len(args) != 1 {
		return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
	}

	return objs.Str(args[0].String()), nil
}

func builtinInt(args []types.Obj) (types.Obj, error) {
	if len(args) != 1 {
		return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
	}

	switch arg := args[0].(type) {
	case objs.Int:
		return arg, nil
	case objs.Str:
		integer, err := strconv.Atoi(string(arg))
		if err != nil {
			return nil, err
		}

		return objs.Int(integer), nil
	case objs.Bool:
		if arg {
			return objs.Int(1), nil
		} else {
			return objs.Int(0), nil
		}
	default:
		return objs.Undefined{}, nil
	}
}
