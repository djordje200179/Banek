package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	errors2 "banek/exec/errors"
	objects2 "banek/exec/objects"
	"banek/interpreter/results"
	stdErrors "errors"
)

func (interpreter *Interpreter) evalExpression(env *environment, expression ast.Expression) (objects2.Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return objects2.Integer(expression), nil
	case expressions.BooleanLiteral:
		return objects2.Boolean(expression), nil
	case expressions.StringLiteral:
		return objects2.String(expression), nil
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

		if condition == objects2.Boolean(true) {
			return interpreter.evalExpression(env, expression.Consequence)
		} else {
			return interpreter.evalExpression(env, expression.Alternative)
		}
	case expressions.Identifier:
		return interpreter.evalIdentifier(env, expression)
	case expressions.FunctionLiteral:
		return objects2.Function{
			Parameters: expression.Parameters,
			Body:       expression.Body,
			Env:        newEnvironment(env),
		}, nil
	case expressions.FunctionCall:
		return interpreter.evalFunctionCall(env, expression)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccess(env, expression)
	default:
		return nil, errors2.UnknownExpressionError{Expression: expression}
	}
}

func (interpreter *Interpreter) evalIdentifier(env *environment, identifier expressions.Identifier) (objects2.Object, error) {
	value, err := env.Get(identifier.String())
	if err == nil {
		return value, nil
	}

	var identifierNotDefinedError errors2.IdentifierNotDefinedError
	if stdErrors.As(err, &identifierNotDefinedError) {
		builtin, ok := objects2.Builtins[identifier.String()]
		if !ok {
			return nil, err
		}

		return builtin, nil
	}

	return nil, err
}

func (interpreter *Interpreter) evalFunctionCall(env *environment, functionCall expressions.FunctionCall) (objects2.Object, error) {
	functionObject, err := interpreter.evalExpression(env, functionCall.Function)
	if err != nil {
		return nil, err
	}

	switch function := functionObject.(type) {
	case objects2.Function:
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
			return objects2.Undefined{}, nil
		}
	case objects2.BuiltinFunction:
		args, err := interpreter.evalFunctionArguments(env, functionCall.Arguments)
		if err != nil {
			return nil, err
		}

		return function(args...)
	default:
		return nil, errors2.InvalidOperandError{Operator: "function call", Operand: functionObject}
	}
}

func (interpreter *Interpreter) evalFunctionArguments(env *environment, expressions []ast.Expression) ([]objects2.Object, error) {
	args := make([]objects2.Object, len(expressions))
	for i, rawArg := range expressions {
		arg, err := interpreter.evalExpression(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

func (interpreter *Interpreter) evalArrayLiteral(env *environment, expression expressions.ArrayLiteral) (objects2.Array, error) {
	elements := make([]objects2.Object, len(expression))
	for i, elementExpression := range expression {
		element, err := interpreter.evalExpression(env, elementExpression)
		if err != nil {
			return nil, err
		}

		elements[i] = element
	}

	return elements, nil
}

func (interpreter *Interpreter) evalCollectionAccess(env *environment, expression expressions.CollectionAccess) (objects2.Object, error) {
	collectionObject, err := interpreter.evalExpression(env, expression.Collection)
	if err != nil {
		return nil, err
	}

	key, err := interpreter.evalExpression(env, expression.Key)
	if err != nil {
		return nil, err
	}

	switch collection := collectionObject.(type) {
	case objects2.Array:
		return interpreter.evalArrayAccess(collection, key)
	default:
		return nil, errors2.InvalidOperandError{Operator: "collection key", Operand: collectionObject}
	}
}

func (interpreter *Interpreter) evalArrayAccess(array objects2.Array, indexObject objects2.Object) (objects2.Object, error) {
	index, ok := indexObject.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: "array index", Operand: indexObject}
	}

	if index < 0 {
		index = objects2.Integer(len(array)) + index
	}

	if index < 0 || index >= objects2.Integer(len(array)) {
		return objects2.Undefined{}, objects2.IndexOutOfBoundsError{Index: int(index), Size: len(array)}
	}

	return array[index], nil
}
