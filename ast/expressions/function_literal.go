package expressions

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type FunctionLiteral struct {
	Parameters []Identifier
	Body       ast.Statement
}

func (literal FunctionLiteral) ExpressionNode() {}

func (literal FunctionLiteral) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.LambdaFunction.String())

	sb.WriteByte('(')
	for i, param := range literal.Parameters {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(literal.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}
