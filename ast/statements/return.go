package statements

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type Return struct {
	Value ast.Expression
}

func (stmt Return) String() string {
	var sb strings.Builder

	sb.WriteString(tokens.Return.String())
	sb.WriteByte(' ')
	sb.WriteString(stmt.Value.String())

	return sb.String()
}

func (stmt Return) HasSideEffects() bool {
	return !stmt.Value.IsConst()
}
