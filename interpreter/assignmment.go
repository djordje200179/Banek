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
		return objs.MakeUndefined(), err
	}

	switch variable := expr.Var.(type) {
	case exprs.Identifier:
		err := env.Set(variable.String(), value)
		if err != nil {
			return objs.MakeUndefined(), err
		}
	case exprs.CollIndex:
		collection, err := interpreter.evalExpr(env, variable.Coll)
		if err != nil {
			return objs.MakeUndefined(), err
		}

		key, err := interpreter.evalExpr(env, variable.Key)
		if err != nil {
			return objs.MakeUndefined(), err
		}

		err = ops.EvalCollSet(collection, key, value)
		if err != nil {
			return objs.MakeUndefined(), err
		}
	default:
		return objs.MakeUndefined(), ast.ErrInvalidAssignment{Variable: expr.Var}
	}

	return value, nil
}
