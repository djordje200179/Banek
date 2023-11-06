package interpreter

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/interpreter/environments"
	"banek/interpreter/results"
)

func (interpreter *interpreter) evalExpr(env environments.Env, expr ast.Expression) (objects.Object, error) {
	switch expr := expr.(type) {
	case expressions.ConstLiteral:
		return expr.Value, nil
	case expressions.ArrayLiteral:
		return interpreter.evalArrayLiteral(env, expr)
	case expressions.UnaryOp:
		return interpreter.evalUnaryOp(env, expr)
	case expressions.BinaryOp:
		return interpreter.evalBinaryOp(env, expr)
	case expressions.Assignment:
		return interpreter.evalAssignment(env, expr)
	case expressions.If:
		cond, err := interpreter.evalExpr(env, expr.Cond)
		if err != nil {
			return nil, err
		}

		if cond == objects.Boolean(true) {
			return interpreter.evalExpr(env, expr.Consequence)
		} else {
			return interpreter.evalExpr(env, expr.Alternative)
		}
	case expressions.Identifier:
		return interpreter.evalIdentifier(env, expr)
	case expressions.FuncLiteral:
		return &environments.Func{
			Params: expr.Params,
			Body:   statements.Return{Value: expr.Body},
			Env:    env,
		}, nil
	case expressions.FuncCall:
		return interpreter.evalFuncCall(env, expr)
	case expressions.CollIndex:
		return interpreter.evalCollIndex(env, expr)
	default:
		return nil, ast.ErrUnknownExpr{Expr: expr}
	}
}

func (interpreter *interpreter) evalBinaryOp(env environments.Env, expr expressions.BinaryOp) (objects.Object, error) {
	right, err := interpreter.evalExpr(env, expr.Right)
	if err != nil {
		return nil, err
	}

	left, err := interpreter.evalExpr(env, expr.Left)
	if err != nil {
		return nil, err
	}

	return operations.EvalBinary(left, right, expr.Operator)
}

func (interpreter *interpreter) evalUnaryOp(env environments.Env, expr expressions.UnaryOp) (objects.Object, error) {
	operand, err := interpreter.evalExpr(env, expr.Operand)
	if err != nil {
		return nil, err
	}

	return operations.EvalUnary(operand, expr.Operation)
}

func (interpreter *interpreter) evalIdentifier(env environments.Env, identifier expressions.Identifier) (objects.Object, error) {
	if index := objects.BuiltinFindIndex(identifier.String()); index != -1 {
		return objects.Builtins[index], nil
	}

	value, err := env.Get(identifier.String())
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (interpreter *interpreter) evalFuncCall(env environments.Env, funcCall expressions.FuncCall) (objects.Object, error) {
	funcObject, err := interpreter.evalExpr(env, funcCall.Func)
	if err != nil {
		return nil, err
	}

	args, err := interpreter.evalFuncArgs(env, funcCall.Args)
	if err != nil {
		return nil, err
	}

	switch function := funcObject.(type) {
	case *environments.Func:
		funcEnv := EnvFactory(function.Env, len(function.Params))
		for i, param := range function.Params {
			err = funcEnv.Define(param.String(), args[i], true)
			if err != nil {
				return nil, err
			}
		}

		switch body := function.Body.(type) {
		case statements.Return:
			return interpreter.evalExpr(funcEnv, body.Value)
		case statements.Block:
			result, err := interpreter.evalBlockStatement(funcEnv, body)
			if err != nil {
				return nil, err
			}

			ret, ok := result.(results.Return)
			if !ok {
				return objects.Undefined{}, nil
			}

			return ret.Value, nil
		default:
			return nil, ast.ErrUnknownStmt{Stmt: body}
		}
	case objects.BuiltinFunc:
		return function.Func(args...)
	default:
		return nil, errors.ErrInvalidOp{Operator: "call", LeftOperand: funcObject}
	}
}

func (interpreter *interpreter) evalFuncArgs(env environments.Env, rawArgs []ast.Expression) ([]objects.Object, error) {
	args := make([]objects.Object, len(rawArgs))
	for i, rawArg := range rawArgs {
		arg, err := interpreter.evalExpr(env, rawArg)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

func (interpreter *interpreter) evalArrayLiteral(env environments.Env, expr expressions.ArrayLiteral) (objects.Array, error) {
	elems := make([]objects.Object, len(expr))
	for i, elemExpr := range expr {
		elem, err := interpreter.evalExpr(env, elemExpr)
		if err != nil {
			return nil, err
		}

		elems[i] = elem
	}

	return elems, nil
}

func (interpreter *interpreter) evalCollIndex(env environments.Env, expr expressions.CollIndex) (objects.Object, error) {
	coll, err := interpreter.evalExpr(env, expr.Coll)
	if err != nil {
		return nil, err
	}

	key, err := interpreter.evalExpr(env, expr.Key)
	if err != nil {
		return nil, err
	}

	return operations.EvalCollGet(coll, key)
}

func (interpreter *interpreter) evalArrayIndex(array objects.Array, indexObject objects.Object) (objects.Object, error) {
	index, ok := indexObject.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: "index", LeftOperand: array, RightOperand: indexObject}
	}

	if index < 0 {
		index = objects.Integer(len(array)) + index
	}

	if index < 0 || index >= objects.Integer(len(array)) {
		return nil, objects.ErrIndexOutOfBounds{Index: int(index), Size: len(array)}
	}

	return array[index], nil
}
