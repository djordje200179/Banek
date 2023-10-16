package statements

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type Return struct {
	Value ast.Expression
}

func (statement Return) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Return.String())
	sb.WriteByte(' ')
	sb.WriteString(statement.Value.String())

	return sb.String()
}

func (statement Return) HasSideEffects() bool {
	return !statement.Value.IsConstant()
}
