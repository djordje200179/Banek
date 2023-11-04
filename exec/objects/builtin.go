package objects

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type BuiltinFunction struct {
	Name     string
	Function func(args ...Object) (Object, error)
}

func (builtin BuiltinFunction) Type() string   { return "builtin" }
func (builtin BuiltinFunction) Clone() Object  { return builtin }
func (builtin BuiltinFunction) String() string { return builtin.Name }

type ErrIncorrectArgumentNumber struct {
	Expected int
	Got      int
}

func (err ErrIncorrectArgumentNumber) Error() string {
	return fmt.Sprintf("incorrect number of arguments: expected %d, got %d", err.Expected, err.Got)
}

var Builtins = []BuiltinFunction{
	{
		Name: "print",
		Function: func(args ...Object) (Object, error) {
			var sb strings.Builder

			for _, arg := range args {
				sb.WriteString(arg.String())
			}

			fmt.Print(sb.String())

			return Undefined{}, nil
		},
	},
	{
		Name: "println",
		Function: func(args ...Object) (Object, error) {
			var sb strings.Builder

			for _, arg := range args {
				sb.WriteString(arg.String())
			}

			fmt.Println(sb.String())

			return Undefined{}, nil
		},
	},
	{
		Name: "read",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 0 {
				return nil, ErrIncorrectArgumentNumber{Expected: 0, Got: len(args)}
			}

			var input string
			_, err := fmt.Scan(&input)
			if err != nil {
				return nil, err
			}

			return String(input), nil
		},
	},
	{
		Name: "readln",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 0 {
				return nil, ErrIncorrectArgumentNumber{Expected: 0, Got: len(args)}
			}

			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return nil, err
			}

			return String(input), nil
		},
	},
	{
		Name: "len",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgumentNumber{Expected: 1, Got: len(args)}
			}

			switch arg := args[0].(type) {
			case String:
				return Integer(len(arg)), nil
			case Array:
				return Integer(len(arg)), nil
			default:
				return Undefined{}, nil
			}
		},
	},
	{
		Name: "str",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgumentNumber{Expected: 1, Got: len(args)}
			}

			return String(args[0].String()), nil
		},
	},
	{
		Name: "int",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgumentNumber{Expected: 1, Got: len(args)}
			}

			switch arg := args[0].(type) {
			case Integer:
				return arg, nil
			case String:
				integer, err := strconv.Atoi(string(arg))
				if err != nil {
					return nil, err
				}

				return Integer(integer), nil
			case Boolean:
				if arg {
					return Integer(1), nil
				} else {
					return Integer(0), nil
				}
			default:
				return Undefined{}, nil
			}
		},
	},
}

func BuiltinFindIndex(name string) int {
	return slices.IndexFunc(Builtins, func(builtin BuiltinFunction) bool {
		return builtin.Name == name
	})
}
