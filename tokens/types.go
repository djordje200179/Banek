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
	PlusAssign
	MinusAssign
	AsteriskAssign
	CaretAssign
	SlashAssign
	ModuloAssign

	Plus
	Minus
	Asterisk
	Slash
	Modulo
	Caret
	Bang

	Arrow

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

var tokenTypeRepresentations = []string{
	Illegal: "illegal",
	EOF:     "EOF",

	Identifier: "identifier",
	Integer:    "integer",
	Boolean:    "boolean",
	String:     "string",
	Undefined:  "undefined",

	Assign:         "=",
	PlusAssign:     "+=",
	MinusAssign:    "-=",
	AsteriskAssign: "*=",
	SlashAssign:    "/=",
	ModuloAssign:   "%=",
	CaretAssign:    "^=",

	Plus:     "+",
	Minus:    "-",
	Asterisk: "*",
	Slash:    "/",
	Modulo:   "%",
	Caret:    "^",
	Bang:     "!",

	Arrow: "->",

	Equals:              "==",
	NotEquals:           "!=",
	LessThan:            "<",
	GreaterThan:         ">",
	LessThanOrEquals:    "<=",
	GreaterThanOrEquals: ">=",

	Comma:     ",",
	SemiColon: ";",

	LeftParenthesis:  "(",
	RightParenthesis: ")",
	LeftBrace:        "{",
	RightBrace:       "}",
	LeftBracket:      "[",
	RightBracket:     "]",

	Function:       "function",
	LambdaFunction: "fn",

	Let:   "let",
	Const: "const",

	Return: "return",

	If:   "if",
	Else: "else",
	Then: "then",

	While: "while",
	Do:    "do",
}

func (tokenType TokenType) String() string {
	return tokenTypeRepresentations[tokenType]
}
