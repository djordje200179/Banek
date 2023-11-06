package statements

import (
	"banek/ast"
	"strings"
)

type If struct {
	Cond ast.Expression

	Consequence, Alternative ast.Statement
}

func (stmt If) String() string {
	var sb strings.Builder

	sb.WriteString("if ")
	sb.WriteString(stmt.Cond.String())
	sb.WriteString(" then {\n")
	sb.WriteString(stmt.Consequence.String())
	sb.WriteString("\n}")
	if stmt.Alternative != nil {
		sb.WriteString(" else {\n")
		sb.WriteString(stmt.Alternative.String())
		sb.WriteString("\n}")
	}

	return sb.String()
}

func (stmt If) HasSideEffects() bool {
	if !stmt.Cond.IsConst() {
		return true
	}

	return stmt.Consequence.HasSideEffects() || stmt.Alternative.HasSideEffects()
}
