package statements

import (
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type FunctionStatement struct {
	Name       expressions.Identifier
	Parameters []expressions.Identifier
	Body       BlockStatement
}

func (function FunctionStatement) StatementNode() {}

func (function FunctionStatement) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Function.String())
	sb.WriteByte(' ')
	sb.WriteString(function.Name.String())

	sb.WriteByte('(')
	for i, param := range function.Parameters {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(function.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}
