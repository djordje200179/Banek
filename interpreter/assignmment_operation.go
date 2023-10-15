package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/interpreter/objects"
)

func (interpreter *Interpreter) evalAssignment(env *environment, expression expressions.InfixOperation, value objects.Object) error {
	switch variable := expression.Left.(type) {
	case expressions.Identifier:
		return interpreter.evalVariableAssignment(env, variable, value)
	case expressions.CollectionAccess:
		return interpreter.evalCollectionAccessAssignment(env, variable, value)
	default:
		return InvalidOperandError{"assignment", nil} // TODO: fix
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
		return InvalidOperandError{"collection key", collectionObject}
	}
}

func (interpreter *Interpreter) evalArrayAccessAssignment(env *environment, array objects.Array, indexExpression ast.Expression, value objects.Object) error {
	indexObject, err := interpreter.evalExpression(env, indexExpression)
	if err != nil {
		return err
	}

	index, ok := indexObject.(objects.Integer)
	if !ok {
		return InvalidOperandError{"array index", indexObject}
	}

	if index < 0 {
		index = objects.Integer(len(array)) + index
	}

	if index < 0 || index >= objects.Integer(len(array)) {
		return IndexOutOfBoundsError{int(index), len(array)}
	}

	array[index] = value

	return nil
}
