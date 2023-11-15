package builtins

import (
	"banek/runtime/types"
)

type BuiltinFunc struct {
	Name    string
	NumArgs int

	Func func(args []types.Obj) (types.Obj, error)
}

func (builtin BuiltinFunc) Type() types.Type { return types.TypeBuiltin }
func (builtin BuiltinFunc) Clone() types.Obj { return builtin }
func (builtin BuiltinFunc) String() string   { return builtin.Name }

func (builtin BuiltinFunc) Equals(other types.Obj) bool {
	otherBuiltin, ok := other.(BuiltinFunc)
	if !ok {
		return false
	}

	return builtin.Name == otherBuiltin.Name
}

var Funcs = [...]BuiltinFunc{
	{
		Name:    "print",
		NumArgs: -1,

		Func: builtinPrint,
	},
	{
		Name:    "println",
		NumArgs: -1,

		Func: builtinPrintln,
	},
	{
		Name:    "read",
		NumArgs: 0,

		Func: builtinRead,
	},
	{
		Name:    "readln",
		NumArgs: 0,

		Func: builtinReadln,
	},
	{
		Name:    "len",
		NumArgs: 1,

		Func: builtinLen,
	},
	{
		Name:    "str",
		NumArgs: 1,

		Func: builtinStr,
	},
	{
		Name:    "int",
		NumArgs: 1,

		Func: builtinInt,
	},
}

func Find(name string) int {
	for i, builtin := range &Funcs {
		if builtin.Name == name {
			return i
		}
	}

	return -1
}
