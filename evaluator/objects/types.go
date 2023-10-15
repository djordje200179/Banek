package objects

type ObjectType int

const (
	NoneType ObjectType = iota

	IntegerType
	BooleanType

	NullType

	ReturnType
)

func (objectType ObjectType) String() string {
	switch objectType {
	case NoneType:
		return "none"

	case IntegerType:
		return "integer"
	case BooleanType:
		return "boolean"

	case NullType:
		return "null"

	case ReturnType:
		return "return"

	default:
		return "UNKNOWN"
	}
}
