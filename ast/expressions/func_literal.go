package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type FuncLiteral struct {
	Params []Identifier
	Body   ast.Expression
}

func (expr FuncLiteral) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Func.String())

	sb.WriteByte('(')
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

func (expr FuncLiteral) IsConst() bool { return false }
