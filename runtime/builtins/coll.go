package builtins

import (
	"banek/runtime/objs"
)

func builtinLen(args []objs.Obj) (objs.Obj, error) {
	res := objs.Obj{
		Type: objs.Int,
	}

	switch args[0].Type {
	case objs.String:
		res.Int = len(args[0].AsString())
	case objs.Array:
		res.Int = len(args[0].AsArray())
	default:
		return objs.Obj{}, InvalidTypeError{
			ArgIndex: 0,
			Arg:      args[0],
		}
	}

	return res, nil
}
