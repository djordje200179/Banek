package interpreter

import (
	"banek/ast/expressions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/interpreter/environments"
	"banek/tokens"
)

var infixOperations = map[tokens.TokenType]operations.InfixOperationType{
	tokens.Plus:     operations.InfixPlusOperation,
	tokens.Minus:    operations.InfixMinusOperation,
	tokens.Asterisk: operations.InfixAsteriskOperation,
	tokens.Slash:    operations.InfixSlashOperation,
	tokens.Modulo:   operations.InfixModuloOperation,
	tokens.Caret:    operations.InfixCaretOperation,

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

func (interpreter *interpreter) evalInfixOperation(env environments.Environment, expression expressions.InfixOperation) (objects.Object, error) {
	right, err := interpreter.evalExpression(env, expression.Right)
	if err != nil {
		return nil, err
	}

	left, err := interpreter.evalExpression(env, expression.Left)
	if err != nil {
		return nil, err
	}

	operation, ok := infixOperations[expression.Operator.Type]
	if !ok {
		return nil, errors.ErrUnknownOperator{Operator: expression.Operator.Type.String()}
	}

	return operations.EvalInfixOperation(left, right, operation)
}

func (interpreter *interpreter) evalPrefixOperation(env environments.Environment, expression expressions.PrefixOperation) (objects.Object, error) {
	operand, err := interpreter.evalExpression(env, expression.Operand)
	if err != nil {
		return nil, err
	}

	operation, ok := prefixOperations[expression.Operator.Type]
	if !ok {
		return nil, errors.ErrUnknownOperator{Operator: expression.Operator.Type.String()}
	}

	return operations.EvalPrefixOperation(operand, operation)
}
