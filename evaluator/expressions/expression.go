package expressions

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/evaluator/objects"
)

func EvalExpression(expression ast.Expression) (objects.Object, error) {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		return objects.Integer(expression.Value), nil
	case expressions.BooleanLiteral:
		return objects.Boolean(expression.Value), nil
	case expressions.PrefixOperation:
		operand, err := EvalExpression(expression.Operand)
		if err != nil {
			return nil, err
		}

		return evalPrefixOperation(expression.Operator, operand)
	case expressions.InfixOperation:
		left, err := EvalExpression(expression.Left)
		if err != nil {
			return nil, err
		}

		right, err := EvalExpression(expression.Right)
		if err != nil {
			return nil, err
		}

		return evalInfixOperation(expression.Operator, left, right)
	case expressions.If:
		condition, err := EvalExpression(expression.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return EvalExpression(expression.Consequence)
		} else {
			return EvalExpression(expression.Alternative)
		}
	default:
		return nil, UnknownExpressionError{expression}
	}
}
