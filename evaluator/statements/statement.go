package statements

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/evaluator/environment"
	"banek/evaluator/expressions"
	"banek/evaluator/objects"
)

func EvalStatement(env *environment.Environment, statement ast.Statement) (objects.Object, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		return expressions.EvalExpression(env, statement.Expression)
	case statements.If:
		condition, err := expressions.EvalExpression(env, statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return EvalStatement(env, statement.Consequence)
		} else if statement.Alternative != nil {
			return EvalStatement(env, statement.Alternative)
		} else {
			return objects.None{}, nil
		}
	case statements.Block:
		blockEnv := environment.New(env)

		for _, statement := range statement.Statements {
			result, err := EvalStatement(blockEnv, statement)
			if err != nil {
				return nil, err
			}

			if result != nil && result.Type() == objects.ReturnType {
				return result, nil
			}
		}

		return objects.None{}, nil
	case statements.Return:
		value, err := expressions.EvalExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		return objects.Return{Value: value}, nil
	case statements.VariableDeclaration:
		value, err := expressions.EvalExpression(env, statement.Value)
		if err != nil {
			return nil, err
		}

		if env.IsDefined(statement.Name.String()) {
			return nil, IdentifierAlreadyDefinedError{statement.Name}
		}

		env.Set(statement.Name.String(), value)

		return objects.None{}, nil
	default:
		return nil, UnknownStatementError{statement}
	}
}
