package exprs

import (
	"banek/ast"
	"banek/symtable"
	"strings"
)

type FuncLiteral struct {
	Params []Ident
	Body   ast.Expr

	Container *symtable.Container
}

func (expr FuncLiteral) String() string {
	var sb strings.Builder

	sb.WriteString("func (")
	for i, param := range expr.Params {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(expr.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (expr FuncLiteral) IsConst() bool {
	return false
}
