package statements

import (
	"banek/ast"
)

type Expr struct {
	Expr ast.Expression
}

func (stmt Expr) String() string {
	return stmt.Expr.String()
}

func (stmt Expr) HasSideEffects() bool {
	return true
}
