package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode"
	"banek/exec/errors"
)

func (compiler *compiler) compileStatement(statement ast.Statement) error {
	switch statement := statement.(type) {
	case statements.Expression:
		err := compiler.compileExpression(statement.Expression)
		if err != nil {
			return err
		}

		compiler.emitInstruction(bytecode.Pop)
	default:
		return errors.UnknownStatementError{Statement: statement}
	}

	return nil
}
