package statements

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/evaluator/expressions"
	"banek/evaluator/objects"
)

func EvalStatement(statement ast.Statement) (objects.Object, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		return expressions.EvalExpression(statement.Expression)
	case statements.If:
		condition, err := expressions.EvalExpression(statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return EvalStatement(statement.Consequence)
		} else if statement.Alternative != nil {
			return EvalStatement(statement.Alternative)
		} else {
			return objects.None{}, nil
		}
	case statements.Block:
		for _, statement := range statement.Statements {
			result, err := EvalStatement(statement)
			if err != nil {
				return nil, err
			}

			if result != nil && result.Type() == objects.ReturnType {
				return result, nil
			}
		}

		return objects.None{}, nil
	case statements.Return:
		value, err := expressions.EvalExpression(statement.Value)
		if err != nil {
			return nil, err
		}

		return objects.Return{Value: value}, nil
	default:
		return nil, UnknownStatementError{statement}
	}
}
