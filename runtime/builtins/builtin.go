package builtins

import (
	"banek/runtime/objs"
	"bytes"
)

type Builtin struct {
	Name    string
	NumArgs int

	Func func(args []objs.Obj) (objs.Obj, error)
}

func GetBuiltin(objs objs.Obj) Builtin {
	return Funcs[objs.IntData]
}

func (builtin Builtin) MakeObj() objs.Obj {
	index := Find(builtin.Name)

	return objs.Obj{Tag: objs.TypeBuiltin, IntData: uint64(index)}
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

func Find(name string) int {
	for i, builtin := range &Funcs {
		if builtin.Name == name {
			return i
		}
	}

	return -1
}

func builtinString(obj objs.Obj) string {
	builtin := GetBuiltin(obj)
	return builtin.Name
}

func builtinEquals(first, second objs.Obj) bool {
	return first.IntData == second.IntData
}

func builtinMarshal(obj objs.Obj) ([]byte, error) {
	builtin := GetBuiltin(obj)

	index := Find(builtin.Name)

	return []byte{byte(index)}, nil
}

func builtinUnmarshal(buf *bytes.Buffer) (objs.Obj, error) {
	index := int(buf.Next(1)[0])

	return Funcs[index].MakeObj(), nil
}

func init() {
	objs.Config[objs.TypeBuiltin] = objs.TypeConfig{
		Stringer: builtinString,
		Equaler:  builtinEquals,

		Marshaller:   builtinMarshal,
		Unmarshaller: builtinUnmarshal,
	}
}
