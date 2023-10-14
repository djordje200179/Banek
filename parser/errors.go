package parser

import "banek/tokens"

type UnexpectedTokenError struct {
	Expected, Got tokens.TokenType
}

func (err UnexpectedTokenError) Error() string {
	return "expected next token to be " + err.Expected.String() + ", got " + err.Got.String() + " instead"
}

type UnknownTokenError struct {
	TokenType tokens.TokenType
}

func (err UnknownTokenError) Error() string {
	return "unknown token type " + err.TokenType.String() + " found"
}
