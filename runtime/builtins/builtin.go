package builtins

import (
	"banek/runtime/types"
)

type BuiltinFunc struct {
	Name string

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

var Builtins = [...]BuiltinFunc{
	{
		Name: "print",
		Func: builtinPrint,
	},
	{
		Name: "println",
		Func: builtinPrintln,
	},
	{
		Name: "read",
		Func: builtinRead,
	},
	{
		Name: "readln",
		Func: builtinReadln,
	},
	{
		Name: "len",
		Func: builtinLen,
	},
	{
		Name: "str",
		Func: builtinStr,
	},
	{
		Name: "int",
		Func: builtinInt,
	},
}
