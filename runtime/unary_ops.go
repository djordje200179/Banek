package runtime

type UnaryOperator uint8

const (
	NegOperator UnaryOperator = iota
	NotOperator
)

func (op UnaryOperator) String() string {
	switch op {
	case NegOperator:
		return "-"
	case NotOperator:
		return "!"
	}

	panic("unreachable")
}
