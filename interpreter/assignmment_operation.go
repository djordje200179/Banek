package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/errors"
	objects2 "banek/exec/objects"
)

func (interpreter *Interpreter) evalAssignment(env *environment, expression expressions.InfixOperation, value objects2.Object) error {
	switch variable := expression.Left.(type) {
	case expressions.Identifier:
		return interpreter.evalVariableAssignment(env, variable, value)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccessAssignment(env, variable, value)
	default:
		return errors.InvalidOperandError{Operator: "assignment"} // TODO: fix
	}
}

func (interpreter *Interpreter) evalVariableAssignment(env *environment, variable expressions.Identifier, value objects2.Object) error {
	return env.Set(variable.String(), value)
}

func (interpreter *Interpreter) evalCollectionAccessAssignment(env *environment, variable expressions.CollectionAccess, value objects2.Object) error {
	collectionObject, err := interpreter.evalExpression(env, variable.Collection)
	if err != nil {
		return err
	}

	switch collection := collectionObject.(type) {
	case objects2.Array:
		return interpreter.evalArrayAccessAssignment(env, collection, variable.Key, value)
	default:
		return errors.InvalidOperandError{Operator: "collection key", Operand: collectionObject}
	}
}

func (interpreter *Interpreter) evalArrayAccessAssignment(env *environment, array objects2.Array, indexExpression ast.Expression, value objects2.Object) error {
	indexObject, err := interpreter.evalExpression(env, indexExpression)
	if err != nil {
		return err
	}

	index, ok := indexObject.(objects2.Integer)
	if !ok {
		return errors.InvalidOperandError{Operator: "array index", Operand: indexObject}
	}

	if index < 0 {
		index = objects2.Integer(len(array)) + index
	}

	if index < 0 || index >= objects2.Integer(len(array)) {
		return objects2.IndexOutOfBoundsError{Index: int(index), Size: len(array)}
	}

	array[index] = value

	return nil
}
