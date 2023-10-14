package tokens

type TokenType int

const (
	Illegal TokenType = iota
	EOF

	Identifier
	Integer
	Boolean

	Assign

	Plus
	Minus
	Asterisk
	Slash

	Bang

	Equals
	NotEquals
	LessThan
	GreaterThan
	LessThanOrEquals
	GreaterThanOrEquals

	Comma
	SemiColon

	LeftParenthesis
	RightParenthesis
	LeftBrace
	RightBrace

	Function
	LambdaFunction

	Var
	Const

	Return
	If
	Else
	While
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case Illegal:
		return "illegal"
	case EOF:
		return "EOF"

	case Identifier:
		return "identifier"
	case Integer:
		return "integer"
	case Boolean:
		return "boolean"

	case Assign:
		return "="

	case Plus:
		return "+"
	case Minus:
		return "-"
	case Asterisk:
		return "*"
	case Slash:
		return "/"

	case Bang:
		return "!"

	case Equals:
		return "=="
	case NotEquals:
		return "!="
	case LessThan:
		return "<"
	case GreaterThan:
		return ">"
	case LessThanOrEquals:
		return "<="
	case GreaterThanOrEquals:
		return ">="

	case Comma:
		return ","
	case SemiColon:
		return ";"

	case LeftParenthesis:
		return "("
	case RightParenthesis:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"

	case Function:
		return "function"
	case LambdaFunction:
		return "fn"

	case Var:
		return "var"
	case Const:
		return "const"

	case Return:
		return "return"
	case If:
		return "if"
	case Else:
		return "else"
	case While:
		return "while"

	default:
		return "UNKNOWN"
	}
}
