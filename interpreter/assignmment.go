package interpreter

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/interpreter/envs"
	"banek/runtime/ops"
	"banek/runtime/types"
)

func (interpreter *interpreter) evalAssignment(env *envs.Env, expr exprs.Assignment) (types.Obj, error) {
	value, err := interpreter.evalExpr(env, expr.Value)
	if err != nil {
		return nil, err
	}

	switch variable := expr.Var.(type) {
	case exprs.Identifier:
		err := env.Set(variable.String(), value)
		if err != nil {
			return nil, err
		}
	case exprs.CollIndex:
		collection, err := interpreter.evalExpr(env, variable.Coll)
		if err != nil {
			return nil, err
		}

		key, err := interpreter.evalExpr(env, variable.Key)
		if err != nil {
			return nil, err
		}

		err = ops.EvalCollSet(collection, key, value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ast.ErrInvalidAssignment{Variable: expr.Var}
	}

	return value, nil
}
