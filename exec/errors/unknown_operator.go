package errors

import "banek/tokens"

type ErrUnknownOperator struct {
	Operator tokens.TokenType
}

func (err ErrUnknownOperator) Error() string {
	return "unknown operator: " + err.Operator.String()
}
