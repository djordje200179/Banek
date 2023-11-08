package objects

type Type byte

const (
	TypeUnknown Type = iota
	TypeUndefined
	TypeBoolean
	TypeInteger
	TypeString
	TypeArray
	TypeFunction
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
