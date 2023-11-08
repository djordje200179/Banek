package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/interpreter/environments"
)

func (interpreter *interpreter) evalAssignment(env environments.Env, expr expressions.Assignment) (objects.Object, error) {
	value, err := interpreter.evalExpr(env, expr.Value)
	if err != nil {
		return nil, err
	}

	switch variable := expr.Var.(type) {
	case expressions.Identifier:
		err := env.Set(variable.String(), value)
		if err != nil {
			return nil, err
		}
	case expressions.CollIndex:
		collection, err := interpreter.evalExpr(env, variable.Coll)
		if err != nil {
			return nil, err
		}

		key, err := interpreter.evalExpr(env, variable.Key)
		if err != nil {
			return nil, err
		}

		err = operations.EvalCollSet(collection, key, value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ast.ErrInvalidAssignment{Variable: expr.Var}
	}

	return value, nil
}
