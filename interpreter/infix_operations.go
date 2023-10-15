package interpreter

import (
	"banek/ast/expressions"
	"banek/interpreter/objects"
	"banek/tokens"
)

type infixOperation func(left, right objects.Object) (objects.Object, error)

var infixOperations = map[tokens.TokenType]infixOperation{
	tokens.Plus:                evalInfixPlusOperation,
	tokens.Minus:               evalInfixMinusOperation,
	tokens.Asterisk:            evalInfixAsteriskOperation,
	tokens.Slash:               evalInfixSlashOperation,
	tokens.Equals:              evalInfixEqualsOperation,
	tokens.NotEquals:           evalInfixNotEqualsOperation,
	tokens.LessThan:            evalInfixLessThanOperation,
	tokens.GreaterThan:         evalInfixGreaterThanOperation,
	tokens.LessThanOrEquals:    evalInfixLessThanOrEqualsOperation,
	tokens.GreaterThanOrEquals: evalInfixGreaterThanOrEqualsOperation,
}

func (interpreter *Interpreter) evalInfixOperation(env *environment, expression expressions.InfixOperation) (objects.Object, error) {
	right, err := interpreter.evalExpression(env, expression.Right)
	if err != nil {
		return nil, err
	}

	switch expression.Operator.Type {
	case tokens.Assign:
		err := interpreter.evalAssignment(env, expression, right)
		if err != nil {
			return nil, err
		}

		return right, nil
	default:
		left, err := interpreter.evalExpression(env, expression.Left)
		if err != nil {
			return nil, err
		}

		operation := infixOperations[expression.Operator.Type]
		if operation == nil {
			return nil, UnknownOperatorError{expression.Operator.Type}
		}

		return operation(left, right)
	}
}

func evalInfixPlusOperation(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return left + right, nil
		default:
			return nil, InvalidOperandError{tokens.Plus.String(), right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return left + right, nil
		default:
			return nil, InvalidOperandError{tokens.Plus.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.Plus.String(), left}
	}
}

func evalInfixMinusOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), right}
	}

	return leftInteger - rightInteger, nil
}

func evalInfixAsteriskOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Asterisk.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Asterisk.String(), right}
	}

	return leftInteger * rightInteger, nil
}

func evalInfixSlashOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Slash.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Slash.String(), right}
	}

	return leftInteger / rightInteger, nil
}

func evalInfixEqualsOperation(left, right objects.Object) (objects.Object, error) {
	return objects.Boolean(left == right), nil
}

func evalInfixNotEqualsOperation(left, right objects.Object) (objects.Object, error) {
	return objects.Boolean(left != right), nil
}

func evalInfixLessThanOperation(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return objects.Boolean(left < right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left < right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.LessThan.String(), left}
	}
}

func evalInfixGreaterThanOperation(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return objects.Boolean(left > right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left > right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.LessThan.String(), left}
	}
}

func evalInfixLessThanOrEqualsOperation(left, right objects.Object) (objects.Object, error) {
	lessThan, err := evalInfixLessThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return lessThan.(objects.Boolean) || equal.(objects.Boolean), nil
}

func evalInfixGreaterThanOrEqualsOperation(left, right objects.Object) (objects.Object, error) {
	greaterThan, err := evalInfixGreaterThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return greaterThan.(objects.Boolean) || equal.(objects.Boolean), nil
}
