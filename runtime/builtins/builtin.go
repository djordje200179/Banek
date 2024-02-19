package builtins

import (
	"banek/runtime"
)

type Builtin struct {
	Name    string
	NumArgs int

	Func func(args []runtime.Obj) (runtime.Obj, error)
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

func (b *Builtin) Truthy() bool       { return true }
func (b *Builtin) String() string     { return b.Name }
func (b *Builtin) Clone() runtime.Obj { return b }

func (b *Builtin) Equals(other runtime.Obj) bool {
	var otherBuiltin *Builtin
	var ok bool

	if otherBuiltin, ok = other.(*Builtin); !ok {
		return false
	}

	return b.Name == otherBuiltin.Name
}

type OperandNotValidError struct {
	Operand runtime.Obj
}

func (e OperandNotValidError) Error() string {
	return "operand type not valid: " + e.Operand.String()
}
