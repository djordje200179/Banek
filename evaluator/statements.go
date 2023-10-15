package evaluator

import (
	"banek/ast"
	"banek/ast/statements"
)

func (evaluator *Evaluator) evaluateStatement(env *environment, statement ast.Statement) (Result, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		value, err := evaluator.evaluateExpression(env, statement.Expression)
		if err != nil {
			return nil, err
		}

		return value, nil
	case statements.If:
		condition, err := evaluator.evaluateExpression(env, statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == Boolean(true) {
			return evaluator.evaluateStatement(env, statement.Consequence)
		} else if statement.Alternative != nil {
			return evaluator.evaluateStatement(env, statement.Alternative)
		} else {
			return None{}, nil
		}
	case statements.Block:
		blockEnv := newEnvironment(env)

		for _, statement := range statement.Statements {
			result, err := evaluator.evaluateStatement(blockEnv, statement)
			if err != nil {
				return nil, err
			}

			switch result := result.(type) {
			case Return:
				return result, nil
			}
		}

		return None{}, nil
	case statements.Return:
		value, err := evaluator.evaluateExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		return Return{Value: value}, nil
	case statements.VariableDeclaration:
		value, err := evaluator.evaluateExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		if env.IsDefined(statement.Name.String()) {
			return nil, IdentifierAlreadyDefinedError{statement.Name}
		}

		env.Set(statement.Name.String(), value)

		return None{}, nil
	case statements.Function:
		value := Function{
			Parameters: statement.Parameters,
			Body:       statement.Body,
			Env:        newEnvironment(env),
		}

		env.Set(statement.Name.String(), value)

		return None{}, nil
	default:
		return nil, UnknownStatementError{statement}
	}
}
