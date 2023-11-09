package exprs

import (
	"banek/ast"
	"strings"
)

type ArrayLiteral []ast.Expr

func (expr ArrayLiteral) String() string {
	var sb strings.Builder

	elemStrings := make([]string, len(expr))
	for i, elem := range expr {
		elemStrings[i] = elem.String()
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elemStrings, ", "))
	sb.WriteByte(']')

	return sb.String()
}

func (expr ArrayLiteral) IsConst() bool {
	for _, elem := range expr {
		if !elem.IsConst() {
			return false
		}
	}

	return true
}
