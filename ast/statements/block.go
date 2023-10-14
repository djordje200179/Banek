package statements

import (
	"banek/ast"
	"strings"
)

type Block struct {
	Statements []ast.Statement
}

func (statement Block) StatementNode() {}

func (statement Block) String() string {
	var sb strings.Builder

	for i, statement := range statement.Statements {
		if i != 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(statement.String())
	}

	return sb.String()
}
