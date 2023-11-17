package builtins

import (
	"banek/runtime/objs"
	"strconv"
)

func builtinStr(args []objs.Obj) (objs.Obj, error) {
	arg := args[0]
	return objs.MakeStr(arg.String()), nil
}

func builtinInt(args []objs.Obj) (objs.Obj, error) {
	arg := args[0]

	switch arg.Tag {
	case objs.TypeInt:
		return arg, nil
	case objs.TypeStr:
		str := arg.AsStr()
		integer, err := strconv.Atoi(str)
		if err != nil {
			return objs.MakeUndefined(), err
		}
		return objs.MakeInt(integer), nil
	case objs.TypeBool:
		boolean := arg.AsBool()
		if boolean {
			return objs.MakeInt(1), nil
		} else {
			return objs.MakeInt(0), nil
		}
	default:
		// TODO: error
		return objs.MakeUndefined(), nil
	}
}
