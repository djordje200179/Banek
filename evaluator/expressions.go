package evaluator

import (
	"banek/ast"
	"banek/ast/expressions"
	"errors"
)

func (evaluator *Evaluator) evaluateExpression(env *environment, expression ast.Expression) (Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return Integer(expression), nil
	case expressions.BooleanLiteral:
		return Boolean(expression), nil
	case expressions.StringLiteral:
		return String(expression), nil
	case expressions.ArrayLiteral:
		elements := make([]Object, len(expression))
		for i, element := range expression {
			evaluatedElement, err := evaluator.evaluateExpression(env, element)
			if err != nil {
				return nil, err
			}

			elements[i] = evaluatedElement
		}

		return Array(elements), nil
	case expressions.VariableAssignment:
		value, err := evaluator.evaluateExpression(env, expression.Value)
		if err != nil {
			return nil, err
		}

		err = env.Set(expression.Variable.String(), value)
		if err != nil {
			return nil, err
		}

		return value, nil
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
		value, err := env.Get(expression.String())
		if err == nil {
			return value, nil
		}

		var identifierNotDefinedError IdentifierNotDefinedError
		if errors.As(err, &identifierNotDefinedError) {
			builtin, ok := builtins[expression.String()]
			if !ok {
				return nil, err
			}

			return builtin, nil
		}

		return nil, err
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

		switch function := functionObject.(type) {
		case Function:
			args, err := evaluator.calculateFunctionArguments(env, expression.Arguments)
			if err != nil {
				return nil, err
			}

			functionEnv := newEnvironment(function.Env)
			for i, param := range function.Parameters {
				err = functionEnv.Define(param.String(), args[i], true)
				if err != nil {
					return nil, err
				}
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
		case BuiltinFunction:
			args, err := evaluator.calculateFunctionArguments(env, expression.Arguments)
			if err != nil {
				return nil, err
			}

			return function(args...)
		default:
			return nil, InvalidOperandError{"function call", functionObject}
		}
	case expressions.CollectionIndex:
		collectionObject, err := evaluator.evaluateExpression(env, expression.Collection)
		if err != nil {
			return nil, err
		}

		switch collection := collectionObject.(type) {
		case Array:
			indexObject, err := evaluator.evaluateExpression(env, expression.Index)
			if err != nil {
				return nil, err
			}

			index, ok := indexObject.(Integer)
			if !ok {
				return nil, InvalidOperandError{"array index", indexObject}
			}

			if index < 0 {
				index = Integer(len(collection)) + index
			}

			if index < 0 || index >= Integer(len(collection)) {
				return Null{}, nil
			}

			return collection[index], nil
		default:
			return nil, InvalidOperandError{"collection index", collectionObject}
		}

	default:
		return nil, UnknownExpressionError{expression}
	}
}

func (evaluator *Evaluator) calculateFunctionArguments(env *environment, rawArgs []ast.Expression) ([]Object, error) {
	args := make([]Object, len(rawArgs))
	for i, rawArg := range rawArgs {
		arg, err := evaluator.evaluateExpression(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}
