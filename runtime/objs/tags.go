package objs

type Tag byte

const (
	TypeUndefined Tag = iota
	TypeBool
	TypeInt
	TypeStr
	TypeArray
	TypeBuiltin
	TypeFunc
)

func (t Tag) String() string {
	return typeNames[t]
}

var typeNames = [...]string{
	"undefined",
	"boolean",
	"integer",
	"string",
	"array",
	"builtin",
	"function",
}
