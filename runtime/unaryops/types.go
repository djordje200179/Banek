package unaryops

type UnaryOperator uint8

const (
	NegOperator UnaryOperator = iota
	NotOperator

	unaryOperatorCount
)

func (op UnaryOperator) String() string {
	switch op {
	case NegOperator:
		return "-"
	case NotOperator:
		return "!"
	default:
		panic("unreachable")
	}
}
