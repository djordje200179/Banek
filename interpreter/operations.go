package interpreter

import (
	"banek/ast/expressions"
	"banek/exec/objects"
	"banek/interpreter/operations"
	"banek/tokens"
)

func (interpreter *Interpreter) evalInfixOperation(env *environment, expression expressions.InfixOperation) (objects.Object, error) {
	right, err := interpreter.evalExpression(env, expression.Right)
	if err != nil {
		return nil, err
	}

	switch expression.Operator.Type {
	case tokens.Assign:
		err := interpreter.evalAssignment(env, expression, right)
		if err != nil {
			return nil, err
		}

		return right, nil
	default:
		left, err := interpreter.evalExpression(env, expression.Left)
		if err != nil {
			return nil, err
		}

		return operations.EvalInfixOperation(left, right, expression.Operator.Type)
	}
}

func (interpreter *Interpreter) evalPrefixOperation(env *environment, expression expressions.PrefixOperation) (objects.Object, error) {
	operand, err := interpreter.evalExpression(env, expression.Operand)
	if err != nil {
		return nil, err
	}

	return operations.EvalPrefixOperation(operand, expression.Operator.Type)
}
