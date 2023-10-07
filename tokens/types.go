package tokens

type TokenType int

const (
	Illegal TokenType = iota
	EOF

	Identifier
	Integer

	Assign

	Plus
	Minus
	Multiply
	Divide

	Bang

	Equals
	NotEquals
	LessThan
	GreaterThan
	LessThanOrEquals
	GreaterThanOrEquals

	Comma
	SemiColon

	LeftParen
	RightParen
	LeftBrace
	RightBrace

	Function
	Let
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case Illegal:
		return "Illegal"
	case EOF:
		return "EOF"

	case Identifier:
		return "Identifier"
	case Integer:
		return "Integer"

	case Assign:
		return "Assign"

	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"

	case Bang:
		return "Bang"

	case Equals:
		return "Equals"
	case NotEquals:
		return "NotEquals"
	case LessThan:
		return "LessThan"
	case GreaterThan:
		return "GreaterThan"
	case LessThanOrEquals:
		return "LessThanOrEquals"

	case Comma:
		return "Comma"
	case SemiColon:
		return "SemiColon"

	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"

	case Function:
		return "Function"
	case Let:
		return "Let"

	default:
		return "UNKNOWN"
	}
}
