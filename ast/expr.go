package ast

import "fmt"

type Expr interface {
	fmt.Stringer

	IsConst() bool
}

type ErrUnknownExpr struct {
	Expr Expr
}

func (err ErrUnknownExpr) Error() string {
	return "unknown expression: " + err.Expr.String()
}

type ErrInvalidAssignment struct {
	Variable Expr
}

func (err ErrInvalidAssignment) Error() string {
	return "invalid assignment to " + err.Variable.String()
}
