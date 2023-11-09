package exprs

import (
	"banek/ast"
	"strings"
)

type FuncCall struct {
	Func ast.Expr
	Args []ast.Expr
}

func (expr FuncCall) String() string {
	var sb strings.Builder

	sb.WriteString(expr.Func.String())
	sb.WriteByte('(')
	for i, arg := range expr.Args {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(arg.String())
	}
	sb.WriteByte(')')

	return sb.String()
}

func (expr FuncCall) IsConst() bool {
	for _, arg := range expr.Args {
		if !arg.IsConst() {
			return false
		}
	}

	return expr.Func.IsConst()
}
