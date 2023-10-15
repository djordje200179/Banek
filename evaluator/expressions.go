package evaluator

import (
	"banek/ast"
	"banek/ast/expressions"
)

func (evaluator *Evaluator) evaluateExpression(env *environment, expression ast.Expression) (Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return Integer(expression.Value), nil
	case expressions.BooleanLiteral:
		return Boolean(expression.Value), nil
	case expressions.PrefixOperation:
		operand, err := evaluator.evaluateExpression(env, expression.Operand)
		if err != nil {
			return nil, err
		}

		return evaluator.evalPrefixOperation(expression.Operator, operand)
	case expressions.InfixOperation:
		left, err := evaluator.evaluateExpression(env, expression.Left)
		if err != nil {
			return nil, err
		}

		right, err := evaluator.evaluateExpression(env, expression.Right)
		if err != nil {
			return nil, err
		}

		return evaluator.evalInfixOperation(expression.Operator, left, right)
	case expressions.If:
		condition, err := evaluator.evaluateExpression(env, expression.Condition)
		if err != nil {
			return nil, err
		}

		if condition == Boolean(true) {
			return evaluator.evaluateExpression(env, expression.Consequence)
		} else {
			return evaluator.evaluateExpression(env, expression.Alternative)
		}
	case expressions.Identifier:
		value, ok := env.Get(expression.String())
		if !ok {
			return nil, IdentifierNotDefinedError{expression}
		}

		return value, nil
	case expressions.FunctionLiteral:
		return Function{
			Parameters: expression.Parameters,
			Body:       expression.Body,
			Env:        newEnvironment(env),
		}, nil
	case expressions.FunctionCall:
		functionObject, err := evaluator.evaluateExpression(env, expression.Function)
		if err != nil {
			return nil, err
		}

		function, ok := functionObject.(Function)
		if !ok {
			return nil, InvalidOperandError{"function call", functionObject}
		}

		functionEnv := newEnvironment(function.Env)
		for i, param := range function.Parameters {
			argExpression := expression.Arguments[i]
			arg, err := evaluator.evaluateExpression(env, argExpression)
			if err != nil {
				return nil, err
			}

			functionEnv.Set(param.Name, arg)
		}

		result, err := evaluator.evaluateStatement(functionEnv, function.Body)
		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case Return:
			return result.Value, nil
		default:
			return Null{}, nil
		}
	default:
		return nil, UnknownExpressionError{expression}
	}
}
