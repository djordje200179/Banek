package builtins

import (
	"banek/runtime/objs"
)

type Builtin struct {
	Name    string
	NumArgs int

	Func func(args []objs.Obj) (objs.Obj, error)
}

var Funcs = [...]Builtin{
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
