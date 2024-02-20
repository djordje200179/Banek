package binaryops

type BinaryOperator uint8

const (
	AddOperator BinaryOperator = iota
	SubOperator
	MulOperator
	DivOperator
	ModOperator

	binaryOperatorCount
)

func (op BinaryOperator) String() string {
	switch op {
	case AddOperator:
		return "+"
	case SubOperator:
		return "-"
	case MulOperator:
		return "*"
	case DivOperator:
		return "/"
	case ModOperator:
		return "%"
	default:
		panic("unreachable")
	}
}
