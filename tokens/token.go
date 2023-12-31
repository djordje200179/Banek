package tokens

import "strings"

type Token struct {
	Type    Type
	Literal string
}

func (token Token) String() string {
	switch token.Type {
	case Identifier, Integer:
		var sb strings.Builder

		sb.WriteString(token.Type.String())
		if token.Literal != "" {
			sb.WriteByte('(')
			sb.WriteString(token.Literal)
			sb.WriteByte(')')
		}

		return sb.String()
	default:
		return token.Type.String()
	}
}
