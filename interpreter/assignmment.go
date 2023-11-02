package interpreter

import (
	"banek/ast/expressions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/interpreter/environments"
)

func (interpreter *interpreter) evalAssignment(env environments.Environment, expression expressions.Assignment) (objects.Object, error) {
	value, err := interpreter.evalExpression(env, expression.Value)
	if err != nil {
		return nil, err
	}

	switch variable := expression.Variable.(type) {
	case expressions.Identifier:
		err := env.Set(variable.String(), value)
		if err != nil {
			return nil, err
		}
	case expressions.CollectionAccess:
		collection, err := interpreter.evalExpression(env, variable.Collection)
		if err != nil {
			return nil, err
		}

		key, err := interpreter.evalExpression(env, variable.Key)
		if err != nil {
			return nil, err
		}

		err = operations.EvalCollectionSet(collection, key, value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidOperand{Operation: "=", LeftOperand: objects.Unknown, RightOperand: value}
	}

	return value, nil
}
