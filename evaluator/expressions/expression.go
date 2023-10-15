package expressions

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/evaluator/environment"
	"banek/evaluator/objects"
)

func EvalExpression(env *environment.Environment, expression ast.Expression) (objects.Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return objects.Integer(expression.Value), nil
	case expressions.BooleanLiteral:
		return objects.Boolean(expression.Value), nil
	case expressions.PrefixOperation:
		operand, err := EvalExpression(env, expression.Operand)
		if err != nil {
			return nil, err
		}

		return evalPrefixOperation(expression.Operator, operand)
	case expressions.InfixOperation:
		left, err := EvalExpression(env, expression.Left)
		if err != nil {
			return nil, err
		}

		right, err := EvalExpression(env, expression.Right)
		if err != nil {
			return nil, err
		}

		return evalInfixOperation(expression.Operator, left, right)
	case expressions.If:
		condition, err := EvalExpression(env, expression.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return EvalExpression(env, expression.Consequence)
		} else {
			return EvalExpression(env, expression.Alternative)
		}
	case expressions.Identifier:
		value, ok := env.Get(expression.String())
		if !ok {
			return nil, IdentifierNotDefinedError{expression}
		}

		return value, nil
	default:
		return nil, UnknownExpressionError{expression}
	}
}
