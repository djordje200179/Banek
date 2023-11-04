package interpreter

import (
	"banek/ast/expressions"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/interpreter/environments"
)

func (interpreter *interpreter) evalInfixOperation(env environments.Environment, expression expressions.InfixOperation) (objects.Object, error) {
	right, err := interpreter.evalExpression(env, expression.Right)
	if err != nil {
		return nil, err
	}

	left, err := interpreter.evalExpression(env, expression.Left)
	if err != nil {
		return nil, err
	}

	return operations.EvalInfixOperation(left, right, expression.Operator)
}

func (interpreter *interpreter) evalPrefixOperation(env environments.Environment, expression expressions.PrefixOperation) (objects.Object, error) {
	operand, err := interpreter.evalExpression(env, expression.Operand)
	if err != nil {
		return nil, err
	}

	return operations.EvalPrefixOperation(operand, expression.Operator)
}
