package statements

import (
	"banek/ast"
	"strings"
)

type Block struct {
	Statements []ast.Statement
}

func (block Block) String() string {
	var sb strings.Builder

	for i, singleStatement := range block.Statements {
		if i != 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(singleStatement.String())
	}

	return sb.String()
}

func (block Block) HasSideEffects() bool {
	for _, singleStatement := range block.Statements {
		if singleStatement.HasSideEffects() {
			return true
		}
	}

	return false
}
