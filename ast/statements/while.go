package statements

import (
	"banek/ast"
	"strings"
)

type While struct {
	Condition ast.Expression

	Body ast.Statement
}

func (statement While) String() string {
	var sb strings.Builder

	sb.WriteString("while")
	sb.WriteString(statement.Condition.String())
	sb.WriteString(" do {\n")
	sb.WriteString(statement.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (statement While) HasSideEffects() bool {
	if !statement.Condition.IsConstant() {
		return true
	}

	return statement.Body.HasSideEffects()
}
