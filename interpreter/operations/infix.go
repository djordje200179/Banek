package operations

import (
	"banek/interpreter/errors"
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

func EvalInfixOperation(left, right objects.Object, operator tokens.TokenType) (objects.Object, error) {
	operation := infixOperations[operator]
	if operation == nil {
		return nil, errors.UnknownOperatorError{Operator: operator}
	}

	return operation(left, right)
}

type InvalidOperandError struct {
}

func evalInfixPlusOperation(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return left + right, nil
		default:
			return nil, errors.InvalidOperandError{Operator: tokens.Plus.String(), Operand: right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return left + right, nil
		default:
			return nil, errors.InvalidOperandError{Operator: tokens.Plus.String(), Operand: right}
		}
	default:
		return nil, errors.InvalidOperandError{Operator: tokens.Plus.String(), Operand: left}
	}
}

func evalInfixMinusOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Minus.String(), Operand: left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Minus.String(), Operand: right}
	}

	return leftInteger - rightInteger, nil
}

func evalInfixAsteriskOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Asterisk.String(), Operand: left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Asterisk.String(), Operand: right}
	}

	return leftInteger * rightInteger, nil
}

func evalInfixSlashOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Slash.String(), Operand: left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Slash.String(), Operand: right}
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
			return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left < right), nil
		default:
			return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	default:
		return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: left}
	}
}

func evalInfixGreaterThanOperation(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return objects.Boolean(left > right), nil
		default:
			return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left > right), nil
		default:
			return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	default:
		return nil, errors.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: left}
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
