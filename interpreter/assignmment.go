package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/interpreter/environments"
)

func (interpreter *interpreter) evalAssignment(env environments.Environment, expression expressions.Assignment) (objects.Object, error) {
	value, err := interpreter.evalExpression(env, expression.Value)
	if err != nil {
		return nil, err
	}

	switch variable := expression.Variable.(type) {
	case expressions.Identifier:
		err := interpreter.evalVariableAssignment(env, variable, value)
		if err != nil {
			return nil, err
		}
	case expressions.CollectionAccess:
		err := interpreter.evalCollectionAccessAssignment(env, variable, value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidOperand{Operation: "=", LeftOperand: objects.Unknown, RightOperand: value}
	}

	return value, nil
}

func (interpreter *interpreter) evalVariableAssignment(env environments.Environment, variable expressions.Identifier, value objects.Object) error {
	return env.Set(variable.String(), value)
}

func (interpreter *interpreter) evalCollectionAccessAssignment(env environments.Environment, variable expressions.CollectionAccess, value objects.Object) error {
	collectionObject, err := interpreter.evalExpression(env, variable.Collection)
	if err != nil {
		return err
	}

	switch collection := collectionObject.(type) {
	case objects.Array:
		err := interpreter.evalArrayAccessAssignment(env, collection, variable.Key, value)
		if err != nil {
			return err
		}
	default:
		return errors.ErrInvalidOperand{Operation: "=", LeftOperand: objects.Unknown, RightOperand: value}
	}

	return nil
}

func (interpreter *interpreter) evalArrayAccessAssignment(env environments.Environment, array objects.Array, indexExpression ast.Expression, value objects.Object) error {
	indexObject, err := interpreter.evalExpression(env, indexExpression)
	if err != nil {
		return err
	}

	index, ok := indexObject.(objects.Integer)
	if !ok {
		return errors.ErrInvalidOperand{Operation: "index", LeftOperand: array, RightOperand: indexObject}
	}

	if index < 0 {
		index = objects.Integer(len(array)) + index
	}

	if index < 0 || index >= objects.Integer(len(array)) {
		return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(array)}
	}

	array[index] = value

	return nil
}
