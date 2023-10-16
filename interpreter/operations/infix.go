package operations

import (
	errors2 "banek/exec/errors"
	objects2 "banek/exec/objects"
	"banek/tokens"
)

type infixOperation func(left, right objects2.Object) (objects2.Object, error)

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

func EvalInfixOperation(left, right objects2.Object, operator tokens.TokenType) (objects2.Object, error) {
	operation := infixOperations[operator]
	if operation == nil {
		return nil, errors2.UnknownOperatorError{Operator: operator}
	}

	return operation(left, right)
}

type InvalidOperandError struct {
}

func evalInfixPlusOperation(left, right objects2.Object) (objects2.Object, error) {
	switch left := left.(type) {
	case objects2.Integer:
		switch right := right.(type) {
		case objects2.Integer:
			return left + right, nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.Plus.String(), Operand: right}
		}
	case objects2.String:
		switch right := right.(type) {
		case objects2.String:
			return left + right, nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.Plus.String(), Operand: right}
		}
	default:
		return nil, errors2.InvalidOperandError{Operator: tokens.Plus.String(), Operand: left}
	}
}

func evalInfixMinusOperation(left, right objects2.Object) (objects2.Object, error) {
	leftInteger, ok := left.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Minus.String(), Operand: left}
	}

	rightInteger, ok := right.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Minus.String(), Operand: right}
	}

	return leftInteger - rightInteger, nil
}

func evalInfixAsteriskOperation(left, right objects2.Object) (objects2.Object, error) {
	leftInteger, ok := left.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Asterisk.String(), Operand: left}
	}

	rightInteger, ok := right.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Asterisk.String(), Operand: right}
	}

	return leftInteger * rightInteger, nil
}

func evalInfixSlashOperation(left, right objects2.Object) (objects2.Object, error) {
	leftInteger, ok := left.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Slash.String(), Operand: left}
	}

	rightInteger, ok := right.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Slash.String(), Operand: right}
	}

	return leftInteger / rightInteger, nil
}

func evalInfixEqualsOperation(left, right objects2.Object) (objects2.Object, error) {
	return objects2.Boolean(left == right), nil
}

func evalInfixNotEqualsOperation(left, right objects2.Object) (objects2.Object, error) {
	return objects2.Boolean(left != right), nil
}

func evalInfixLessThanOperation(left, right objects2.Object) (objects2.Object, error) {
	switch left := left.(type) {
	case objects2.Integer:
		switch right := right.(type) {
		case objects2.Integer:
			return objects2.Boolean(left < right), nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	case objects2.String:
		switch right := right.(type) {
		case objects2.String:
			return objects2.Boolean(left < right), nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	default:
		return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: left}
	}
}

func evalInfixGreaterThanOperation(left, right objects2.Object) (objects2.Object, error) {
	switch left := left.(type) {
	case objects2.Integer:
		switch right := right.(type) {
		case objects2.Integer:
			return objects2.Boolean(left > right), nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	case objects2.String:
		switch right := right.(type) {
		case objects2.String:
			return objects2.Boolean(left > right), nil
		default:
			return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: right}
		}
	default:
		return nil, errors2.InvalidOperandError{Operator: tokens.LessThan.String(), Operand: left}
	}
}

func evalInfixLessThanOrEqualsOperation(left, right objects2.Object) (objects2.Object, error) {
	lessThan, err := evalInfixLessThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return lessThan.(objects2.Boolean) || equal.(objects2.Boolean), nil
}

func evalInfixGreaterThanOrEqualsOperation(left, right objects2.Object) (objects2.Object, error) {
	greaterThan, err := evalInfixGreaterThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return greaterThan.(objects2.Boolean) || equal.(objects2.Boolean), nil
}
