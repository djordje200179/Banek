package exprs

import (
	"banek/ast"
	"strings"
)

type If struct {
	Cond ast.Expr

	Consequence, Alternative ast.Expr
}

func (expr If) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(expr.Cond.String())
	sb.WriteString(" then {\n")
	sb.WriteString(expr.Consequence.String())
	sb.WriteString("\n}")
	if expr.Alternative != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(expr.Alternative.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}

func (expr If) IsConst() bool {
	return expr.Cond.IsConst() && expr.Consequence.IsConst() && expr.Alternative.IsConst()
}
