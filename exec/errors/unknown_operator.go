package errors

import "banek/tokens"

type UnknownOperatorError struct {
	Operator tokens.TokenType
}

func (err UnknownOperatorError) Error() string {
	return "unknown operator: " + err.Operator.String()
}
