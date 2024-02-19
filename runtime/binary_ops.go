package runtime

type BinaryOperator uint8

const (
	AddOperator BinaryOperator = iota
	SubOperator
	MulOperator
	DivOperator
	ModOperator

	LtOperator
	LtEqOperator
	GtOperator
	GtEqOperator
	EqOperator
	NotEqOperator
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

	case LtOperator:
		return "<"
	case LtEqOperator:
		return "<="
	case GtOperator:
		return ">"
	case GtEqOperator:
		return ">="
	case EqOperator:
		return "=="
	case NotEqOperator:
		return "!="
	}

	panic("unreachable")
}
