package compiler

import (
	"banek/ast/expressions"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/operations"
	"banek/tokens"
)

var infixOperations = map[tokens.TokenType]operations.InfixOperationType{
	tokens.Plus:     operations.InfixPlusOperation,
	tokens.Minus:    operations.InfixMinusOperation,
	tokens.Asterisk: operations.InfixAsteriskOperation,
	tokens.Slash:    operations.InfixSlashOperation,

	tokens.Equals:              operations.InfixEqualsOperation,
	tokens.NotEquals:           operations.InfixNotEqualsOperation,
	tokens.LessThan:            operations.InfixLessThanOperation,
	tokens.GreaterThan:         operations.InfixGreaterThanOperation,
	tokens.LessThanOrEquals:    operations.InfixLessThanOrEqualsOperation,
	tokens.GreaterThanOrEquals: operations.InfixGreaterThanOrEqualsOperation,
}

var prefixOperations = map[tokens.TokenType]operations.PrefixOperationType{
	tokens.Minus: operations.PrefixMinusOperation,
	tokens.Bang:  operations.PrefixBangOperation,
}

func (compiler *compiler) compileInfixOperation(expression expressions.InfixOperation) error {
	err := compiler.compileExpression(expression.Left)
	if err != nil {
		return err
	}

	err = compiler.compileExpression(expression.Right)
	if err != nil {
		return err
	}

	operator := expression.Operator.Type
	operation, ok := infixOperations[operator]
	if !ok {
		return errors.ErrUnknownOperator{Operator: operator.String()}
	}

	container := compiler.topContainer()
	container.emitInstruction(instruction.OperationInfix, int(operation))

	return nil
}

func (compiler *compiler) compilePrefixOperation(expression expressions.PrefixOperation) error {
	err := compiler.compileExpression(expression.Operand)
	if err != nil {
		return err
	}

	operator := expression.Operator.Type
	operation, ok := prefixOperations[operator]
	if !ok {
		return errors.ErrUnknownOperator{Operator: operator.String()}
	}

	container := compiler.topContainer()
	container.emitInstruction(instruction.OperationPrefix, int(operation))

	return nil
}
