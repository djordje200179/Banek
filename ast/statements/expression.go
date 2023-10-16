package statements

import (
	"banek/ast"
)

type Expression struct {
	Expression ast.Expression
}

func (statement Expression) String() string {
	return statement.Expression.String()
}

func (statement Expression) HasSideEffects() bool {
	return true
}
