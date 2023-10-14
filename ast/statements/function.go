package statements

import (
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type Function struct {
	Name       expressions.Identifier
	Parameters []expressions.Identifier
	Body       Block
}

func (statement Function) StatementNode() {}

func (statement Function) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Function.String())
	sb.WriteByte(' ')
	sb.WriteString(statement.Name.String())

	sb.WriteByte('(')
	for i, param := range statement.Parameters {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(statement.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}
