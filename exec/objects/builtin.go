package objects

import (
	"fmt"
	"strconv"
	"strings"
)

type BuiltinFunc struct {
	Name string

	Func func(args ...Object) (Object, error)
}

func (builtin BuiltinFunc) Type() Type     { return TypeBuiltin }
func (builtin BuiltinFunc) Clone() Object  { return builtin }
func (builtin BuiltinFunc) String() string { return builtin.Name }

type ErrIncorrectArgNum struct {
	Expected int
	Got      int
}

func (err ErrIncorrectArgNum) Error() string {
	var sb strings.Builder

	sb.WriteString("incorrect number of arguments: expected ")
	sb.WriteString(strconv.Itoa(err.Expected))
	sb.WriteString(", got ")
	sb.WriteString(strconv.Itoa(err.Got))

	return sb.String()
}

var Builtins = [...]BuiltinFunc{
	{
		Name: "print",
		Func: func(args ...Object) (Object, error) {
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
		Func: func(args ...Object) (Object, error) {
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
		Func: func(args ...Object) (Object, error) {
			if len(args) != 0 {
				return nil, ErrIncorrectArgNum{Expected: 0, Got: len(args)}
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
		Func: func(args ...Object) (Object, error) {
			if len(args) != 0 {
				return nil, ErrIncorrectArgNum{Expected: 0, Got: len(args)}
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
		Func: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
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
		Func: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
			}

			return String(args[0].String()), nil
		},
	},
	{
		Name: "int",
		Func: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, ErrIncorrectArgNum{Expected: 1, Got: len(args)}
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
	for i, builtin := range &Builtins {
		if builtin.Name == name {
			return i
		}
	}

	return -1
}
