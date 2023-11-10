package types

type Type byte

const (
	TypeUnknown Type = iota
	TypeUndefined
	TypeBool
	TypeInt
	TypeStr
	TypeArray
	TypeFunc
	TypeBuiltin
)

func (t Type) String() string {
	return typeNames[t]
}

var typeNames = [...]string{
	"unknown",
	"undefined",
	"boolean",
	"integer",
	"string",
	"array",
	"function",
	"builtin",
}