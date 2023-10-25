package errors

import "banek/ast"

type ErrUnknownExpression struct {
	Expression ast.Expression
}

func (err ErrUnknownExpression) Error() string {
	return "unknown expression: " + err.Expression.String()
}

type ErrUnknownStatement struct {
	Statement ast.Statement
}

func (err ErrUnknownStatement) Error() string {
	return "unknown statement: " + err.Statement.String()
}
