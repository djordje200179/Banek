package statements

import (
	"banek/ast"
	"strings"
)

type Block struct {
	Stmts []ast.Statement
}

func (stmt Block) String() string {
	var sb strings.Builder

	for i, stmt := range stmt.Stmts {
		if i != 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(stmt.String())
	}

	return sb.String()
}

func (stmt Block) HasSideEffects() bool {
	for _, stmt := range stmt.Stmts {
		if stmt.HasSideEffects() {
			return true
		}
	}

	return false
}
