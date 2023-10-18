package errors

import "banek/ast"

type ErrUnknownExpression struct {
	Expression ast.Expression
}

func (err ErrUnknownExpression) Error() string {
	return "unknown expression: " + err.Expression.String()
}
