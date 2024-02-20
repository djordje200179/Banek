package builtins

import (
	"banek/runtime/objs"
	"strconv"
)

func builtinStr(args []objs.Obj) (objs.Obj, error) {
	return objs.MakeString(args[0].String()), nil
}

func builtinInt(args []objs.Obj) (objs.Obj, error) {
	res := objs.Obj{Type: objs.Int}

	switch args[0].Type {
	case objs.Int:
		res.Int = args[0].Int
	case objs.String:
		val, err := strconv.Atoi(args[0].AsString())
		if err != nil {
			return objs.Obj{}, err
		}

		res.Int = val
	case objs.Bool:
		res.Int = args[0].Int
	default:
		return objs.Obj{}, InvalidTypeError{
			ArgIndex: 0,
			Arg:      args[0],
		}
	}

	return res, nil
}
