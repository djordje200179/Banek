package statements

import (
	"banek/ast"
	"strings"
)

type BlockStatement struct {
	Statements []ast.Statement
}

func (block BlockStatement) StatementNode() {}

func (block BlockStatement) String() string {
	var sb strings.Builder

	for i, statement := range block.Statements {
		if i != 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(statement.String())
	}

	return sb.String()
}
