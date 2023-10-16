package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/exec/errors"
)

func (compiler *compiler) compileStatement(statement ast.Statement) error {
	switch statement := statement.(type) {
	case statements.Expression:
		return compiler.compileExpression(statement.Expression)
	default:
		return errors.UnknownStatementError{Statement: statement}
	}
}
