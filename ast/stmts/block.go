package stmts

import (
	"banek/ast"
	"strings"
)

type Block []ast.Stmt

func (stmt Block) String() string {
	var sb strings.Builder

	for i, stmt := range stmt {
		if i != 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(stmt.String())
	}

	return sb.String()
}

func (stmt Block) StmtNode() {}
