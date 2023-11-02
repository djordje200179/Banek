package ast

import "fmt"

type Expression interface {
	fmt.Stringer

	IsConstant() bool
}

type ErrUnknownExpression struct {
	Expression Expression
}

func (err ErrUnknownExpression) Error() string {
	return "unknown expression: " + err.Expression.String()
}
