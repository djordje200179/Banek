package statements

import (
	"banek/ast"
	"banek/ast/expressions"
)

type UnknownStatementError struct {
	Statement ast.Statement
}

func (err UnknownStatementError) Error() string {
	return "unknown statement: " + err.Statement.String()
}

type IdentifierAlreadyDefinedError struct {
	Identifier expressions.Identifier
}

func (err IdentifierAlreadyDefinedError) Error() string {
	return "identifier already defined: " + err.Identifier.String()
}
