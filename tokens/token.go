package tokens

import "strings"

type Type byte

const (
	Illegal Type = iota
	EOF

	Ident
	Int
	Bool
	String
	Undefined

	Assign
	PlusAssign
	MinusAssign
	AsteriskAssign
	SlashAssign
	PercentAssign

	Plus
	Minus
	Asterisk
	Slash
	Percent
	Bang

	Equals
	NotEquals
	Less
	Greater
	LessEquals
	GreaterEquals

	Comma
	SemiColon

	LParen
	RParen
	LBrace
	RBrace
	LBracket
	RBracket

	Func
	Return

	LArrow
	RArrow
	VBar

	Let
	Mut

	If
	Else
	Then

	While
	Do
	For

	tokenCount
)

var names = [tokenCount]string{
	Illegal: "ILLEGAL",
	EOF:     "EOF",

	Ident:     "IDENT",
	Int:       "INT",
	Bool:      "BOOL",
	String:    "STRING",
	Undefined: "UNDEFINED",

	Assign:         "=",
	PlusAssign:     "+=",
	MinusAssign:    "-=",
	AsteriskAssign: "*=",
	SlashAssign:    "/=",
	PercentAssign:  "%=",

	Plus:     "+",
	Minus:    "-",
	Asterisk: "*",
	Slash:    "/",
	Percent:  "%",
	Bang:     "!",

	Equals:        "==",
	NotEquals:     "!=",
	Less:          "<",
	Greater:       ">",
	LessEquals:    "<=",
	GreaterEquals: ">=",

	Comma:     ",",
	SemiColon: ";",

	LParen:   "(",
	RParen:   ")",
	LBrace:   "{",
	RBrace:   "}",
	LBracket: "[",
}

func init() {
	for keyword, tokenType := range keywords {
		names[tokenType] = keyword
	}
}

func (t Type) String() string {
	if t >= tokenCount {
		return "UNKNOWN"
	}

	return names[t]
}

type Token struct {
	Type    Type
	Literal string
}

func (t Token) String() string {
	var sb strings.Builder

	sb.WriteString(t.Type.String())
	if t.Literal != "" {
		sb.WriteByte('(')
		sb.WriteString(t.Literal)
		sb.WriteByte(')')
	}

	return sb.String()
}
