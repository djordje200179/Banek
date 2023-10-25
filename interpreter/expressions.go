package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/interpreter/results"
)

func (interpreter *interpreter) evalExpression(env *environment, expression ast.Expression) (objects.Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return objects.Integer(expression), nil
	case expressions.BooleanLiteral:
		return objects.Boolean(expression), nil
	case expressions.StringLiteral:
		return objects.String(expression), nil
	case expressions.ArrayLiteral:
		return interpreter.evalArrayLiteral(env, expression)
	case expressions.PrefixOperation:
		return interpreter.evalPrefixOperation(env, expression)
	case expressions.InfixOperation:
		return interpreter.evalInfixOperation(env, expression)
	case expressions.Assignment:
		return interpreter.evalAssignment(env, expression)
	case expressions.If:
		condition, err := interpreter.evalExpression(env, expression.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return interpreter.evalExpression(env, expression.Consequence)
		} else {
			return interpreter.evalExpression(env, expression.Alternative)
		}
	case expressions.Identifier:
		return interpreter.evalIdentifier(env, expression)
	case expressions.FunctionLiteral:
		return objects.Function{
			Parameters: expression.Parameters,
			Body:       expression.Body,
			Env:        env,
		}, nil
	case expressions.FunctionCall:
		return interpreter.evalFunctionCall(env, expression)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccess(env, expression)
	default:
		return nil, errors.ErrUnknownExpression{Expression: expression}
	}
}

func (interpreter *interpreter) evalIdentifier(env *environment, identifier expressions.Identifier) (objects.Object, error) {
	builtin, ok := objects.Builtins[identifier.String()]
	if ok {
		return builtin, nil
	}

	value, err := env.Get(identifier.String())
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (interpreter *interpreter) evalFunctionCall(env *environment, functionCall expressions.FunctionCall) (objects.Object, error) {
	functionObject, err := interpreter.evalExpression(env, functionCall.Function)
	if err != nil {
		return nil, err
	}

	switch function := functionObject.(type) {
	case objects.Function:
		args, err := interpreter.evalFunctionArguments(env, functionCall.Arguments)
		if err != nil {
			return nil, err
		}

		functionEnv := newEnvironment(function.Env.(*environment), len(function.Parameters))
		for i, param := range function.Parameters {
			err = functionEnv.Define(param.String(), args[i], true)
			if err != nil {
				return nil, err
			}
		}

		switch body := function.Body.(type) {
		case statements.Expression:
			return interpreter.evalExpression(functionEnv, body.Expression)
		case statements.Block:
			result, err := interpreter.evalBlockStatement(functionEnv, body)
			if err != nil {
				return nil, err
			}

			returnValue, ok := result.(results.Return)
			if !ok {
				return objects.Undefined{}, nil
			}

			return returnValue.Value, nil
		default:
			return nil, errors.ErrUnknownStatement{Statement: body}
		}
	case objects.BuiltinFunction:
		args, err := interpreter.evalFunctionArguments(env, functionCall.Arguments)
		if err != nil {
			return nil, err
		}

		return function(args...)
	default:
		return nil, errors.ErrInvalidOperand{Operation: "call", LeftOperand: functionObject}
	}
}

func (interpreter *interpreter) evalFunctionArguments(env *environment, expressions []ast.Expression) ([]objects.Object, error) {
	args := make([]objects.Object, len(expressions))
	for i, rawArg := range expressions {
		arg, err := interpreter.evalExpression(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

func (interpreter *interpreter) evalArrayLiteral(env *environment, expression expressions.ArrayLiteral) (objects.Array, error) {
	elements := make([]objects.Object, len(expression))
	for i, elementExpression := range expression {
		element, err := interpreter.evalExpression(env, elementExpression)
		if err != nil {
			return nil, err
		}

		elements[i] = element
	}

	return elements, nil
}

func (interpreter *interpreter) evalCollectionAccess(env *environment, expression expressions.CollectionAccess) (objects.Object, error) {
	collectionObject, err := interpreter.evalExpression(env, expression.Collection)
	if err != nil {
		return nil, err
	}

	key, err := interpreter.evalExpression(env, expression.Key)
	if err != nil {
		return nil, err
	}

	switch collection := collectionObject.(type) {
	case objects.Array:
		return interpreter.evalArrayAccess(collection, key)
	default:
		return nil, errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: key}
	}
}

func (interpreter *interpreter) evalArrayAccess(array objects.Array, indexObject objects.Object) (objects.Object, error) {
	index, ok := indexObject.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: "index", LeftOperand: array, RightOperand: indexObject}
	}

	if index < 0 {
		index = objects.Integer(len(array)) + index
	}

	if index < 0 || index >= objects.Integer(len(array)) {
		return objects.Undefined{}, objects.ErrIndexOutOfBounds{Index: int(index), Size: len(array)}
	}

	return array[index], nil
}
