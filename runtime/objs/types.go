package objs

type Type byte

const (
	Undefined Type = iota
	Int
	Bool
	String
	Array
	Func
	Builtin

	TypeCount
)

const TypeBits = 3
const TypeMask = (1 << TypeBits) - 1
const IntMask = ^TypeMask
