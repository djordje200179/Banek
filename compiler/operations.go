package compiler

import (
	"banek/ast/expressions"
	"banek/bytecode/instruction"
)

func (compiler *compiler) compileInfixOperation(expression expressions.InfixOperation) error {
	err := compiler.compileExpression(expression.Left)
	if err != nil {
		return err
	}

	err = compiler.compileExpression(expression.Right)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instruction.OperationInfix, int(expression.Operator))

	return nil
}

func (compiler *compiler) compilePrefixOperation(expression expressions.PrefixOperation) error {
	err := compiler.compileExpression(expression.Operand)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instruction.OperationPrefix, int(expression.Operator))

	return nil
}
