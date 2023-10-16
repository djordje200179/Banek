package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/interpreter/errors"
	"banek/interpreter/objects"
	"banek/interpreter/results"
	stdErrors "errors"
)

func (interpreter *Interpreter) evalExpression(env *environment, expression ast.Expression) (objects.Object, error) {
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
			Env:        newEnvironment(env),
		}, nil
	case expressions.FunctionCall:
		return interpreter.evalFunctionCall(env, expression)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccess(env, expression)
	default:
		return nil, errors.UnknownExpressionError{Expression: expression}
	}
}

func (interpreter *Interpreter) evalIdentifier(env *environment, identifier expressions.Identifier) (objects.Object, error) {
	value, err := env.Get(identifier.String())
	if err == nil {
		return value, nil
	}

	var identifierNotDefinedError errors.IdentifierNotDefinedError
	if stdErrors.As(err, &identifierNotDefinedError) {
		builtin, ok := objects.Builtins[identifier.String()]
		if !ok {
			return nil, err
		}

		return builtin, nil
	}

	return nil, err
}

func (interpreter *Interpreter) evalFunctionCall(env *environment, functionCall expressions.FunctionCall) (objects.Object, error) {
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

		functionEnv := newEnvironment(function.Env.(*environment))
		for i, param := range function.Parameters {
			err = functionEnv.Define(param.String(), args[i], true)
			if err != nil {
				return nil, err
			}
		}

		result, err := interpreter.evalStatement(functionEnv, function.Body)
		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case results.Return:
			return result.Value, nil
		default:
			return objects.Undefined{}, nil
		}
	case objects.BuiltinFunction:
		args, err := interpreter.evalFunctionArguments(env, functionCall.Arguments)
		if err != nil {
			return nil, err
		}

		return function(args...)
	default:
		return nil, errors.InvalidOperandError{Operator: "function call", Operand: functionObject}
	}
}

func (interpreter *Interpreter) evalFunctionArguments(env *environment, expressions []ast.Expression) ([]objects.Object, error) {
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

func (interpreter *Interpreter) evalArrayLiteral(env *environment, expression expressions.ArrayLiteral) (objects.Array, error) {
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

func (interpreter *Interpreter) evalCollectionAccess(env *environment, expression expressions.CollectionAccess) (objects.Object, error) {
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
		return nil, errors.InvalidOperandError{Operator: "collection key", Operand: collectionObject}
	}
}

func (interpreter *Interpreter) evalArrayAccess(array objects.Array, indexObject objects.Object) (objects.Object, error) {
	index, ok := indexObject.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: "array index", Operand: indexObject}
	}

	if index < 0 {
		index = objects.Integer(len(array)) + index
	}

	if index < 0 || index >= objects.Integer(len(array)) {
		return objects.Undefined{}, objects.IndexOutOfBoundsError{Index: int(index), Size: len(array)}
	}

	return array[index], nil
}
