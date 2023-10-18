package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/errors"
	"banek/exec/objects"
)

func (interpreter *Interpreter) evalAssignment(env *environment, expression expressions.InfixOperation, value objects.Object) error {
	switch variable := expression.Left.(type) {
	case expressions.Identifier:
		return interpreter.evalVariableAssignment(env, variable, value)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccessAssignment(env, variable, value)
	default:
		return errors.ErrInvalidOperand{Operator: "assignment"} // TODO: fix
	}
}

func (interpreter *Interpreter) evalVariableAssignment(env *environment, variable expressions.Identifier, value objects.Object) error {
	return env.Set(variable.String(), value)
}

func (interpreter *Interpreter) evalCollectionAccessAssignment(env *environment, variable expressions.CollectionAccess, value objects.Object) error {
	collectionObject, err := interpreter.evalExpression(env, variable.Collection)
	if err != nil {
		return err
	}

	switch collection := collectionObject.(type) {
	case objects.Array:
		return interpreter.evalArrayAccessAssignment(env, collection, variable.Key, value)
	default:
		return errors.ErrInvalidOperand{Operator: "collection key", Operand: collectionObject}
	}
}

func (interpreter *Interpreter) evalArrayAccessAssignment(env *environment, array objects.Array, indexExpression ast.Expression, value objects.Object) error {
	indexObject, err := interpreter.evalExpression(env, indexExpression)
	if err != nil {
		return err
	}

	index, ok := indexObject.(objects.Integer)
	if !ok {
		return errors.ErrInvalidOperand{Operator: "array index", Operand: indexObject}
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
