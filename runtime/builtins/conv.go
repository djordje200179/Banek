package builtins

import (
	"banek/runtime/objs"
	"strconv"
)

func builtinStr(args []objs.Obj) (objs.Obj, error) {
	return objs.MakeString(args[0].String()), nil
}

func builtinInt(args []objs.Obj) (objs.Obj, error) {
	switch args[0].Type() {
	case objs.Int:
		return objs.MakeInt(args[0].AsInt()), nil
	case objs.String:
		val, err := strconv.Atoi(args[0].AsString())
		if err != nil {
			return objs.Obj{}, err
		}

		return objs.MakeInt(val), nil
	case objs.Bool:
		if args[0].AsBool() {
			return objs.MakeInt(1), nil
		} else {
			return objs.MakeInt(0), nil
		}
	default:
		return objs.Obj{}, InvalidTypeError{
			ArgIndex: 0,
			Arg:      args[0],
		}
	}
}
