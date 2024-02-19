package exprs

import (
	"banek/ast"
	"strings"
)

type If struct {
	Cond ast.Expr

	Cons, Alt ast.Expr
}

func (expr If) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(expr.Cond.String())
	sb.WriteString(" then {\n")
	sb.WriteString(expr.Cons.String())
	sb.WriteString("\n}")
	if expr.Alt != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(expr.Alt.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}

func (expr If) IsConst() bool {
	return expr.Cond.IsConst() && expr.Cons.IsConst() && expr.Alt.IsConst()
}
