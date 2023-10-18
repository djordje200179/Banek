package parser

import "banek/tokens"

type ErrUnexpectedToken struct {
	Expected, Got tokens.TokenType
}

func (err ErrUnexpectedToken) Error() string {
	return "expected next token to be " + err.Expected.String() + ", got " + err.Got.String() + " instead"
}

type ErrUnknownToken struct {
	TokenType tokens.TokenType
}

func (err ErrUnknownToken) Error() string {
	return "unknown token type " + err.TokenType.String() + " found"
}
