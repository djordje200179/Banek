package interpreter

import (
	"banek/ast"
	"banek/ast/stmts"
	"banek/interpreter/envs"
	"banek/interpreter/results"
	"banek/runtime/objs"
)

func (interpreter *interpreter) evalStmt(env *envs.Env, stmt ast.Stmt) (results.Result, error) {
	switch stmt := stmt.(type) {
	case stmts.Expr:
		return interpreter.evalExpr(env, stmt.Expr)
	case stmts.If:
		cond, err := interpreter.evalExpr(env, stmt.Cond)
		if err != nil {
			return nil, err
		}

		if cond == objs.Bool(true) {
			return interpreter.evalStmt(env, stmt.Consequence)
		} else if stmt.Alternative != nil {
			return interpreter.evalStmt(env, stmt.Alternative)
		} else {
			return results.None{}, nil
		}
	case stmts.While:
		for {
			cond, err := interpreter.evalExpr(env, stmt.Cond)
			if err != nil {
				return nil, err
			}

			if cond != objs.Bool(true) {
				break
			}

			result, err := interpreter.evalStmt(env, stmt.Body)
			if err != nil {
				return nil, err
			}

			switch result := result.(type) {
			case results.Return:
				return result, nil
			}
		}

		return results.None{}, nil
	case stmts.Block:
		blockEnv := envs.New(env, 0)
		return interpreter.evalBlock(blockEnv, stmt)
	case stmts.Return:
		value, err := interpreter.evalExpr(env, stmt.Value)
		if err != nil {
			return nil, err
		}

		return results.Return{Value: value}, nil
	case stmts.VarDecl:
		value, err := interpreter.evalExpr(env, stmt.Value)
		if err != nil {
			return nil, err
		}

		err = env.Define(stmt.Name.String(), value, stmt.Mutable)
		if err != nil {
			return nil, err
		}

		return results.None{}, nil
	case stmts.Func:
		value := &envs.Func{
			Params: stmt.Params,
			Body:   stmt.Body,
			Env:    env,
		}

		err := env.Define(stmt.Name.String(), value, false)
		if err != nil {
			return nil, err
		}

		return results.None{}, nil
	case stmts.Invalid:
		return nil, results.Error{Err: stmt.Err}
	default:
		return nil, ast.ErrUnknownStmt{Stmt: stmt}
	}
}

func (interpreter *interpreter) evalBlock(env *envs.Env, block stmts.Block) (results.Result, error) {
	for _, stmt := range block.Stmts {
		result, err := interpreter.evalStmt(env, stmt)
		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case results.Return:
			return result, nil
		}
	}

	return results.None{}, nil
}
