package exprs

import (
	"banek/ast"
	"strings"
)

type CollIndex struct {
	Coll, Key ast.Expr
}

func (expr CollIndex) String() string {
	var sb strings.Builder

	sb.WriteString(expr.Coll.String())
	sb.WriteByte('[')
	sb.WriteString(expr.Key.String())
	sb.WriteByte(']')

	return sb.String()
}

func (expr CollIndex) IsConst() bool {
	return expr.Coll.IsConst() && expr.Key.IsConst()
}
