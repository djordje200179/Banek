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
	case statements.If:
		err := compiler.compileExpression(statement.Condition)
		if err != nil {
			return err
		}

		firstPatchAddress := compiler.currentAddress()
		compiler.emitInstruction(bytecode.BranchIfFalse, 0)

		err = compiler.compileStatement(statement.Consequence)
		if err != nil {
			return err
		}

		elseAddress := compiler.currentAddress()

		if statement.Alternative != nil {
			secondPatchAddress := compiler.currentAddress()
			compiler.emitInstruction(bytecode.Branch, 0)
			elseAddress += bytecode.Branch.Info().Size()

			err = compiler.compileStatement(statement.Alternative)
			if err != nil {
				return err
			}

			outAddress := compiler.currentAddress()
			compiler.patchInstructionOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-bytecode.Branch.Info().Size())
		}

		compiler.patchInstructionOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-bytecode.BranchIfFalse.Info().Size())

	case statements.Block:
		for _, statement := range statement.Statements {
			err := compiler.compileStatement(statement)
			if err != nil {
				return err
			}
		}
	default:
		return errors.UnknownStatementError{Statement: statement}
	}

	return nil
}
