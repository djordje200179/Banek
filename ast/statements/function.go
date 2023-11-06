package statements

import (
	"banek/ast/expressions"
	"banek/tokens"
	"strings"
)

type Func struct {
	Name expressions.Identifier

	Params []expressions.Identifier
	Body   Block
}

func (stmt Func) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Func.String())
	sb.WriteByte(' ')
	sb.WriteString(stmt.Name.String())

	sb.WriteByte('(')
	for i, param := range stmt.Params {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(stmt.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (stmt Func) HasSideEffects() bool {
	return true
}
