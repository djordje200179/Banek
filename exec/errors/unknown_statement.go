package errors

import "banek/ast"

type ErrUnknownStatement struct {
	Statement ast.Statement
}

func (err ErrUnknownStatement) Error() string {
	return "unknown statement: " + err.Statement.String()
}
