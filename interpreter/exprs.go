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
	"banek/runtime/types"
)

func (interpreter *interpreter) evalExpr(env *envs.Env, expr ast.Expr) (types.Obj, error) {
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
			return nil, err
		}

		if cond == objs.Bool(true) {
			return interpreter.evalExpr(env, expr.Consequence)
		} else {
			return interpreter.evalExpr(env, expr.Alternative)
		}
	case exprs.Identifier:
		return interpreter.evalIdentifier(env, expr)
	case exprs.FuncLiteral:
		return &envs.Func{
			Params: expr.Params,
			Body:   stmts.Return{Value: expr.Body},
			Env:    env,
		}, nil
	case exprs.FuncCall:
		return interpreter.evalFuncCall(env, expr)
	case exprs.CollIndex:
		return interpreter.evalCollIndex(env, expr)
	default:
		return nil, ast.ErrUnknownExpr{Expr: expr}
	}
}

func (interpreter *interpreter) evalBinaryOp(env *envs.Env, expr exprs.BinaryOp) (types.Obj, error) {
	right, err := interpreter.evalExpr(env, expr.Right)
	if err != nil {
		return nil, err
	}

	left, err := interpreter.evalExpr(env, expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator >= ops.BinaryOperator(len(ops.BinaryOps)) {
		return nil, errors.ErrUnknownOperator{Operator: expr.Operator.String()}
	}

	return ops.BinaryOps[expr.Operator](left, right)
}

func (interpreter *interpreter) evalUnaryOp(env *envs.Env, expr exprs.UnaryOp) (types.Obj, error) {
	operand, err := interpreter.evalExpr(env, expr.Operand)
	if err != nil {
		return nil, err
	}

	if expr.Operator >= ops.UnaryOperator(len(ops.UnaryOps)) {
		return nil, errors.ErrUnknownOperator{Operator: expr.Operator.String()}
	}

	return ops.UnaryOps[expr.Operator](operand)
}

func (interpreter *interpreter) evalIdentifier(env *envs.Env, identifier exprs.Identifier) (types.Obj, error) {
	if index := builtins.BuiltinFindIndex(identifier.String()); index != -1 {
		return builtins.Funcs[index], nil
	}

	value, err := env.Get(identifier.String())
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (interpreter *interpreter) evalFuncCall(env *envs.Env, funcCall exprs.FuncCall) (types.Obj, error) {
	funcObject, err := interpreter.evalExpr(env, funcCall.Func)
	if err != nil {
		return nil, err
	}

	args, err := interpreter.evalFuncArgs(env, funcCall.Args)
	if err != nil {
		return nil, err
	}

	switch function := funcObject.(type) {
	case *envs.Func:
		funcEnv := envs.New(function.Env, len(function.Params))
		for i, param := range function.Params {
			err = funcEnv.Define(param.String(), args[i], true)
			if err != nil {
				return nil, err
			}
		}

		switch body := function.Body.(type) {
		case stmts.Return:
			return interpreter.evalExpr(funcEnv, body.Value)
		case stmts.Block:
			result, err := interpreter.evalBlock(funcEnv, body)
			if err != nil {
				return nil, err
			}

			ret, ok := result.(results.Return)
			if !ok {
				return objs.Undefined{}, nil
			}

			return ret.Value, nil
		default:
			return nil, ast.ErrUnknownStmt{Stmt: body}
		}
	case builtins.BuiltinFunc:
		return function.Func(args)
	default:
		return nil, errors.ErrInvalidOp{Operator: "call", LeftOperand: funcObject}
	}
}

func (interpreter *interpreter) evalFuncArgs(env *envs.Env, rawArgs []ast.Expr) ([]types.Obj, error) {
	args := make([]types.Obj, len(rawArgs))
	for i, rawArg := range rawArgs {
		arg, err := interpreter.evalExpr(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

func (interpreter *interpreter) evalArrayLiteral(env *envs.Env, expr exprs.ArrayLiteral) (*objs.Array, error) {
	array := &objs.Array{
		Slice: make([]types.Obj, len(expr)),
	}

	for i, elemExpr := range expr {
		elem, err := interpreter.evalExpr(env, elemExpr)
		if err != nil {
			return nil, err
		}

		array.Slice[i] = elem
	}

	return array, nil
}

func (interpreter *interpreter) evalCollIndex(env *envs.Env, expr exprs.CollIndex) (types.Obj, error) {
	coll, err := interpreter.evalExpr(env, expr.Coll)
	if err != nil {
		return nil, err
	}

	key, err := interpreter.evalExpr(env, expr.Key)
	if err != nil {
		return nil, err
	}

	return ops.EvalCollGet(coll, key)
}
