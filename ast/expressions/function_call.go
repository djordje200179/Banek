package expressions

import (
	"banek/ast"
	"strings"
)

type FunctionCall struct {
	Function  ast.Expression
	Arguments []ast.Expression
}

func (call FunctionCall) String() string {
	var sb strings.Builder

	sb.WriteString(call.Function.String())
	sb.WriteByte('(')
	for i, arg := range call.Arguments {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(arg.String())
	}
	sb.WriteByte(')')

	return sb.String()
}

func (call FunctionCall) IsConstant() bool {
	for _, arg := range call.Arguments {
		if !arg.IsConstant() {
			return false
		}
	}

	return call.Function.IsConstant()
}
