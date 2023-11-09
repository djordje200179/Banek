package exprs

import "banek/ast"

type Assignment struct {
	Var, Value ast.Expr
}

func (expr Assignment) String() string {
	return expr.Var.String() + " = " + expr.Value.String()
}

func (expr Assignment) IsConst() bool {
	return expr.Value.IsConst()
}
