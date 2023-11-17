package interpreter

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/ast/stmts"
	"banek/interpreter/envs"
	"banek/interpreter/results"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/ops"
)

func (interpreter *interpreter) evalExpr(env *envs.Env, expr ast.Expr) (objs.Obj, error) {
	switch expr := expr.(type) {
	case exprs.ConstLiteral:
		return expr.Value, nil
	case exprs.ArrayLiteral:
		return interpreter.evalArrayLiteral(env, expr)
	case exprs.UnaryOp:
		return interpreter.evalUnaryOp(env, expr)
	case exprs.BinaryOp:
		return interpreter.evalBinaryOp(env, expr)
	case exprs.Assignment:
		return interpreter.evalAssignment(env, expr)
	case exprs.If:
		cond, err := interpreter.evalExpr(env, expr.Cond)
		if err != nil {
			return objs.Obj{}, err
		}

		if cond.Tag != objs.TypeBool {
			return objs.Obj{}, errors.ErrNotBool{Obj: cond}
		}

		if cond.AsBool() {
			return interpreter.evalExpr(env, expr.Consequence)
		} else {
			return interpreter.evalExpr(env, expr.Alternative)
		}
	case exprs.Identifier:
		return interpreter.evalIdentifier(env, expr)
	case exprs.FuncLiteral:
		function := &envs.Func{
			Params: expr.Params,
			Body:   stmts.Return{Value: expr.Body},
			Env:    env,
		}

		return function.MakeObj(), nil
	case exprs.FuncCall:
		return interpreter.evalFuncCall(env, expr)
	case exprs.CollIndex:
		return interpreter.evalCollIndex(env, expr)
	default:
		return objs.Obj{}, ast.ErrUnknownExpr{Expr: expr}
	}
}

func (interpreter *interpreter) evalBinaryOp(env *envs.Env, expr exprs.BinaryOp) (objs.Obj, error) {
	right, err := interpreter.evalExpr(env, expr.Right)
	if err != nil {
		return objs.Obj{}, err
	}

	left, err := interpreter.evalExpr(env, expr.Left)
	if err != nil {
		return objs.Obj{}, err
	}

	if expr.Operator >= ops.BinaryOperator(len(ops.BinaryOps)) {
		return objs.Obj{}, errors.ErrUnknownOperator{Operator: expr.Operator.String()}
	}

	return ops.BinaryOps[expr.Operator](left, right)
}

func (interpreter *interpreter) evalUnaryOp(env *envs.Env, expr exprs.UnaryOp) (objs.Obj, error) {
	operand, err := interpreter.evalExpr(env, expr.Operand)
	if err != nil {
		return objs.Obj{}, err
	}

	if expr.Operator >= ops.UnaryOperator(len(ops.UnaryOps)) {
		return objs.Obj{}, errors.ErrUnknownOperator{Operator: expr.Operator.String()}
	}

	return ops.UnaryOps[expr.Operator](operand)
}

func (interpreter *interpreter) evalIdentifier(env *envs.Env, identifier exprs.Identifier) (objs.Obj, error) {
	if index := builtins.Find(identifier.String()); index != -1 {
		builtin := &builtins.Funcs[index]
		return builtin.MakeObj(), nil
	}

	value, err := env.Get(identifier.String())
	if err != nil {
		return objs.Obj{}, err
	}

	return value, nil
}

func (interpreter *interpreter) evalFuncCall(env *envs.Env, funcCall exprs.FuncCall) (objs.Obj, error) {
	funcObj, err := interpreter.evalExpr(env, funcCall.Func)
	if err != nil {
		return objs.Obj{}, err
	}

	args, err := interpreter.evalFuncArgs(env, funcCall.Args)
	if err != nil {
		return objs.Obj{}, err
	}

	switch funcObj.Tag {
	case objs.TypeFunc:
		function := envs.GetFunc(funcObj)

		if len(args) > len(function.Params) {
			return objs.Obj{}, errors.ErrTooManyArgs{Expected: len(function.Params), Received: len(args)}
		}

		funcEnv := envs.New(function.Env, len(function.Params))
		for i, param := range function.Params {
			err = funcEnv.Define(param.String(), args[i], true)
			if err != nil {
				return objs.Obj{}, err
			}
		}

		switch body := function.Body.(type) {
		case stmts.Return:
			return interpreter.evalExpr(funcEnv, body.Value)
		case stmts.Block:
			result, err := interpreter.evalBlock(funcEnv, body)
			if err != nil {
				return objs.Obj{}, err
			}

			ret, ok := result.(results.Return)
			if !ok {
				return objs.Obj{}, nil
			}

			return ret.Value, nil
		default:
			return objs.Obj{}, ast.ErrUnknownStmt{Stmt: body}
		}
	case objs.TypeBuiltin:
		builtin := builtins.GetBuiltin(funcObj)

		if builtin.NumArgs != -1 && len(args) != builtin.NumArgs {
			return objs.Obj{}, errors.ErrTooManyArgs{Expected: builtin.NumArgs, Received: len(args)}
		}

		return builtin.Func(args)
	default:
		return objs.Obj{}, errors.ErrNotCallable{Obj: funcObj}
	}
}

func (interpreter *interpreter) evalFuncArgs(env *envs.Env, rawArgs []ast.Expr) ([]objs.Obj, error) {
	args := make([]objs.Obj, len(rawArgs))
	for i, rawArg := range rawArgs {
		arg, err := interpreter.evalExpr(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

func (interpreter *interpreter) evalArrayLiteral(env *envs.Env, expr exprs.ArrayLiteral) (objs.Obj, error) {
	array := &objs.Array{
		Slice: make([]objs.Obj, len(expr)),
	}

	for i, elemExpr := range expr {
		elem, err := interpreter.evalExpr(env, elemExpr)
		if err != nil {
			return objs.Obj{}, err
		}

		array.Slice[i] = elem
	}

	return objs.MakeArray(array), nil
}

func (interpreter *interpreter) evalCollIndex(env *envs.Env, expr exprs.CollIndex) (objs.Obj, error) {
	coll, err := interpreter.evalExpr(env, expr.Coll)
	if err != nil {
		return objs.Obj{}, err
	}

	key, err := interpreter.evalExpr(env, expr.Key)
	if err != nil {
		return objs.Obj{}, err
	}

	return ops.EvalCollGet(coll, key)
}
