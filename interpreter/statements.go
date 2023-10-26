package interpreter

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/exec/environments"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/interpreter/results"
)

func (interpreter *interpreter) evalStatement(env environments.Environment, statement ast.Statement) (results.Result, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		return interpreter.evalExpression(env, statement.Expression)
	case statements.If:
		condition, err := interpreter.evalExpression(env, statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return interpreter.evalStatement(env, statement.Consequence)
		} else if statement.Alternative != nil {
			return interpreter.evalStatement(env, statement.Alternative)
		} else {
			return results.None, nil
		}
	case statements.While:
		for {
			condition, err := interpreter.evalExpression(env, statement.Condition)
			if err != nil {
				return nil, err
			}

			if condition != objects.Boolean(true) {
				break
			}

			result, err := interpreter.evalStatement(env, statement.Body)
			if err != nil {
				return nil, err
			}

			switch result := result.(type) {
			case results.Return:
				return result, nil
			}
		}

		return results.None, nil
	case statements.Block:
		blockEnv := EnvFactory(env, 0)
		return interpreter.evalBlockStatement(blockEnv, statement)
	case statements.Return:
		value, err := interpreter.evalExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		return results.Return{Value: value}, nil
	case statements.VariableDeclaration:
		value, err := interpreter.evalExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		err = env.Define(statement.Name.String(), value, !statement.Const)
		if err != nil {
			return nil, err
		}

		return results.None, nil
	case statements.Function:
		value := &objects.Function{
			Parameters: statement.Parameters,
			Body:       statement.Body,
			Env:        env,
		}

		err := env.Define(statement.Name.String(), value, false)
		if err != nil {
			return nil, err
		}

		return results.None, nil
	case statements.Error:
		return nil, results.Error{Err: statement.Err}
	default:
		return nil, errors.ErrUnknownStatement{Statement: statement}
	}
}

func (interpreter *interpreter) evalBlockStatement(env environments.Environment, block statements.Block) (results.Result, error) {
	for _, statement := range block.Statements {
		result, err := interpreter.evalStatement(env, statement)
		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case results.Return:
			return result, nil
		}
	}

	return results.None, nil
}
