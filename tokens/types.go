package tokens

type TokenType int

const (
	Illegal TokenType = iota
	EOF

	Identifier
	Integer
	Boolean
	String
	Undefined

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
	LeftBracket
	RightBracket

	Function
	LambdaFunction

	Let
	Const

	Return

	If
	Else
	Then

	While
	Do
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
	case String:
		return "string"
	case Undefined:
		return "undefined"

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
	case LeftBracket:
		return "["
	case RightBracket:
		return "]"

	case Function:
		return "function"
	case LambdaFunction:
		return "fn"

	case Let:
		return "let"
	case Const:
		return "const"

	case Return:
		return "return"

	case If:
		return "if"
	case Else:
		return "else"
	case Then:
		return "then"

	case While:
		return "while"
	case Do:
		return "do"

	default:
		return "UNKNOWN"
	}
}
