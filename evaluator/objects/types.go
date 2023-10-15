package objects

type ObjectType int

const (
	NoneType ObjectType = iota

	IntegerType
	BooleanType

	NullType
)

func (objectType ObjectType) String() string {
	switch objectType {
	case IntegerType:
		return "integer"
	case BooleanType:
		return "boolean"
	case NullType:
		return "null"
	case NoneType:
		return "none"

	default:
		return "UNKNOWN"
	}
}
