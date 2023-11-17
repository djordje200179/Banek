package interpreter

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/interpreter/envs"
	"banek/runtime/objs"
	"banek/runtime/ops"
)

func (interpreter *interpreter) evalAssignment(env *envs.Env, expr exprs.Assignment) (objs.Obj, error) {
	value, err := interpreter.evalExpr(env, expr.Value)
	if err != nil {
		return objs.Obj{}, err
	}

	switch variable := expr.Var.(type) {
	case exprs.Identifier:
		err := env.Set(variable.String(), value)
		if err != nil {
			return objs.Obj{}, err
		}
	case exprs.CollIndex:
		collection, err := interpreter.evalExpr(env, variable.Coll)
		if err != nil {
			return objs.Obj{}, err
		}

		key, err := interpreter.evalExpr(env, variable.Key)
		if err != nil {
			return objs.Obj{}, err
		}

		err = ops.EvalCollSet(collection, key, value)
		if err != nil {
			return objs.Obj{}, err
		}
	default:
		return objs.Obj{}, ast.ErrInvalidAssignment{Variable: expr.Var}
	}

	return value, nil
}
