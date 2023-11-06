package parser

import (
	"banek/tokens"
	"strings"
)

type ErrUnexpectedToken struct {
	Expected, Got tokens.Type
}

func (err ErrUnexpectedToken) Error() string {
	var sb strings.Builder

	sb.WriteString("expected ")
	sb.WriteString(err.Expected.String())
	sb.WriteString(", got ")
	sb.WriteString(err.Got.String())
	sb.WriteString(" instead")

	return sb.String()
}

type ErrUnknownToken struct {
	TokenType tokens.Type
}

func (err ErrUnknownToken) Error() string {
	var sb strings.Builder

	sb.WriteString("unknown token type ")
	sb.WriteString(err.TokenType.String())

	return sb.String()
}
