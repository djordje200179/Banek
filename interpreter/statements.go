package interpreter

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/exec/errors"
	objects2 "banek/exec/objects"
	"banek/interpreter/results"
)

func (interpreter *Interpreter) evalStatement(env *environment, statement ast.Statement) (Result, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		value, err := interpreter.evalExpression(env, statement.Expression)
		if err != nil {
			return nil, err
		}

		return value, nil
	case statements.If:
		condition, err := interpreter.evalExpression(env, statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects2.Boolean(true) {
			return interpreter.evalStatement(env, statement.Consequence)
		} else if statement.Alternative != nil {
			return interpreter.evalStatement(env, statement.Alternative)
		} else {
			return results.None{}, nil
		}
	case statements.Block:
		blockEnv := newEnvironment(env)

		for _, statement := range statement.Statements {
			result, err := interpreter.evalStatement(blockEnv, statement)
			if err != nil {
				return nil, err
			}

			switch result := result.(type) {
			case results.Return:
				return result, nil
			}
		}

		return results.None{}, nil
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

		return results.None{}, nil
	case statements.Function:
		value := objects2.Function{
			Parameters: statement.Parameters,
			Body:       statement.Body,
			Env:        newEnvironment(env),
		}

		err := env.Define(statement.Name.String(), value, false)
		if err != nil {
			return nil, err
		}

		return results.None{}, nil
	case statements.Error:
		return nil, results.Error{Err: statement.Err}
	default:
		return nil, errors.UnknownStatementError{Statement: statement}
	}
}
