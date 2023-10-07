package tokens

import "strings"

type Token struct {
	Type    TokenType
	Literal string
}

func (token Token) String() string {
	switch token.Type {
	case Identifier, Integer:
		var sb strings.Builder

		sb.WriteString(token.Type.String())
		sb.WriteRune('(')
		sb.WriteString(token.Literal)
		sb.WriteRune(')')

		return sb.String()
	default:
		return token.Type.String()
	}
}
