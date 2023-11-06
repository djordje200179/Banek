package tokens

type Type byte

const (
	Illegal Type = iota
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

	Equals
	NotEquals
	Less
	Greater
	LessEquals
	GreaterEquals

	Comma
	SemiColon

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Func
	Return

	VerticalBar
	Arrow

	Let
	Mut

	If
	Else
	Then

	While
	Do
)

var typeStrings = [...]string{
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

	Equals:        "==",
	NotEquals:     "!=",
	Less:          "<",
	Greater:       ">",
	LessEquals:    "<=",
	GreaterEquals: ">=",

	Comma:     ",",
	SemiColon: ";",

	LeftParen:    "(",
	RightParen:   ")",
	LeftBrace:    "{",
	RightBrace:   "}",
	LeftBracket:  "[",
	RightBracket: "]",

	Func:   "func",
	Return: "return",

	Arrow:       "->",
	VerticalBar: "|",

	Let: "let",
	Mut: "mut",

	If:   "if",
	Else: "else",
	Then: "then",

	While: "while",
	Do:    "do",
}

func (t Type) String() string {
	return typeStrings[t]
}
