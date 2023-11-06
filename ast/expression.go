package ast

import "fmt"

type Expression interface {
	fmt.Stringer

	IsConst() bool
}

type ErrUnknownExpr struct {
	Expr Expression
}

func (err ErrUnknownExpr) Error() string {
	return "unknown expression: " + err.Expr.String()
}
