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
