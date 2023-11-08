package interpreter

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/exec/objects"
	"banek/interpreter/environments"
	"banek/interpreter/results"
)

func (interpreter *interpreter) evalStmt(env *environments.Env, stmt ast.Statement) (results.Result, error) {
	switch stmt := stmt.(type) {
	case statements.Expr:
		return interpreter.evalExpr(env, stmt.Expr)
	case statements.If:
		cond, err := interpreter.evalExpr(env, stmt.Cond)
		if err != nil {
			return nil, err
		}

		if cond == objects.Boolean(true) {
			return interpreter.evalStmt(env, stmt.Consequence)
		} else if stmt.Alternative != nil {
			return interpreter.evalStmt(env, stmt.Alternative)
		} else {
			return results.None{}, nil
		}
	case statements.While:
		for {
			cond, err := interpreter.evalExpr(env, stmt.Cond)
			if err != nil {
				return nil, err
			}

			if cond != objects.Boolean(true) {
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
	case statements.Block:
		blockEnv := environments.New(env, 0)
		return interpreter.evalBlockStatement(blockEnv, stmt)
	case statements.Return:
		value, err := interpreter.evalExpr(env, stmt.Value)
		if err != nil {
			return nil, err
		}

		return results.Return{Value: value}, nil
	case statements.VarDecl:
		value, err := interpreter.evalExpr(env, stmt.Value)
		if err != nil {
			return nil, err
		}

		err = env.Define(stmt.Name.String(), value, stmt.Mutable)
		if err != nil {
			return nil, err
		}

		return results.None{}, nil
	case statements.Func:
		value := &environments.Func{
			Params: stmt.Params,
			Body:   stmt.Body,
			Env:    env,
		}

		err := env.Define(stmt.Name.String(), value, false)
		if err != nil {
			return nil, err
		}

		return results.None{}, nil
	case statements.Invalid:
		return nil, results.Error{Err: stmt.Err}
	default:
		return nil, ast.ErrUnknownStmt{Stmt: stmt}
	}
}

func (interpreter *interpreter) evalBlockStatement(env *environments.Env, block statements.Block) (results.Result, error) {
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
