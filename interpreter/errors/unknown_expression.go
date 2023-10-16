package errors

import "banek/ast"

type UnknownExpressionError struct {
	Expression ast.Expression
}

func (err UnknownExpressionError) Error() string {
	return "unknown expression: " + err.Expression.String()
}
