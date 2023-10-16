package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode/instruction"
	"banek/exec/errors"
)

func (compiler *compiler) compileStatement(statement ast.Statement) error {
	switch statement := statement.(type) {
	case statements.Expression:
		err := compiler.compileExpression(statement.Expression)
		if err != nil {
			return err
		}

		compiler.emitInstruction(instruction.Pop)
	case statements.If:
		err := compiler.compileExpression(statement.Condition)
		if err != nil {
			return err
		}

		firstPatchAddress := compiler.currentAddress()
		compiler.emitInstruction(instruction.BranchIfFalse, 0)

		err = compiler.compileStatement(statement.Consequence)
		if err != nil {
			return err
		}

		elseAddress := compiler.currentAddress()

		if statement.Alternative != nil {
			branchSize := instruction.Branch.Info().Size()

			secondPatchAddress := elseAddress
			compiler.emitInstruction(instruction.Branch, 0)
			elseAddress += branchSize

			err = compiler.compileStatement(statement.Alternative)
			if err != nil {
				return err
			}

			outAddress := compiler.currentAddress()
			compiler.patchInstructionOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-branchSize)
		}

		compiler.patchInstructionOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

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
