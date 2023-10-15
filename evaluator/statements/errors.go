package statements

import "banek/ast"

type UnknownStatementError struct {
	Statement ast.Statement
}

func (err UnknownStatementError) Error() string {
	return "unknown statement: " + err.Statement.String()
}
