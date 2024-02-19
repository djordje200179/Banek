package builtins

import (
	"banek/runtime"
	"banek/runtime/primitives"
	"strconv"
)

func builtinStr(args []runtime.Obj) (runtime.Obj, error) {
	return primitives.String(args[0].String()), nil
}

func builtinInt(args []runtime.Obj) (runtime.Obj, error) {
	switch arg := args[0].(type) {
	case primitives.Int:
		return arg, nil
	case primitives.String:
		val, err := strconv.Atoi(string(arg))
		if err != nil {
			return nil, err
		}

		return primitives.Int(val), nil
	case primitives.Bool:
		if arg {
			return primitives.Int(1), nil
		} else {
			return primitives.Int(0), nil
		}
	default:
		return nil, runtime.InvalidTypeError{
			BuiltinName: "int",
			ArgIndex:    0,
			Arg:         args[0],
		}
	}
}
