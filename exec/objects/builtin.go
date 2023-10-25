package objects

import (
	"fmt"
	"strings"
)

type BuiltinFunction func(args ...Object) (Object, error)

func (builtin BuiltinFunction) Type() string   { return "builtin" }
func (builtin BuiltinFunction) String() string { return "builtin function" }

type ErrIncorrectArgumentNumber struct {
	Expected int
	Got      int
}

func (err ErrIncorrectArgumentNumber) Error() string {
	return fmt.Sprintf("incorrect number of arguments: expected %d, got %d", err.Expected, err.Got)
}

var Builtins = map[string]BuiltinFunction{
	"len": func(args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, ErrIncorrectArgumentNumber{Expected: 1, Got: len(args)}
		}

		switch arg := args[0].(type) {
		case String:
			return Integer(len(arg)), nil
		case Array:
			return Integer(len(arg)), nil
		default:
			return Undefined, nil
		}
	},
	"print": func(args ...Object) (Object, error) {
		var sb strings.Builder

		for _, arg := range args {
			sb.WriteString(arg.String())
		}

		fmt.Println(sb.String())

		return Undefined, nil
	},
	"str": func(args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, ErrIncorrectArgumentNumber{Expected: 1, Got: len(args)}
		}

		return String(args[0].String()), nil
	},
}
