package builtins

import (
	"banek/runtime/objs"
)

func builtinLen(args []objs.Obj) (objs.Obj, error) {
	switch args[0].Type() {
	case objs.String:
		return objs.MakeInt(len(args[0].AsString())), nil
	case objs.Array:
		return objs.MakeInt(len(args[0].AsArray())), nil
	default:
		return objs.Obj{}, InvalidTypeError{
			ArgIndex: 0,
			Arg:      args[0],
		}
	}
}
