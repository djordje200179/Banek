package statements

import (
	"banek/ast"
)

type ExpressionStatement struct {
	Expression ast.Expression
}

func (statement ExpressionStatement) StatementNode() {}

func (statement ExpressionStatement) String() string {
	return statement.Expression.String()
}
