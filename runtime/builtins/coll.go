package builtins

import (
	"banek/runtime/objs"
)

func builtinLen(args []objs.Obj) (objs.Obj, error) {
	arg := args[0]
	switch arg.Tag {
	case objs.TypeStr:
		str := arg.AsStr()
		return objs.MakeInt(len(str)), nil
	case objs.TypeArray:
		arr := arg.AsArray()
		return objs.MakeInt(len(arr.Slice)), nil
	default:
		// TODO: error
		return objs.Obj{}, nil
	}
}
