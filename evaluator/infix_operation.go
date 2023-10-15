package evaluator

import (
	"banek/tokens"
)

type infixOperation func(left, right Object) (Object, error)

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

func (evaluator *Evaluator) evalInfixOperation(operator tokens.Token, left, right Object) (Object, error) {
	operation := infixOperations[operator.Type]
	if operation == nil {
		return nil, UnknownOperatorError{operator.Type}
	}

	return operation(left, right)
}

func evalInfixPlusOperation(left, right Object) (Object, error) {
	switch left := left.(type) {
	case Integer:
		switch right := right.(type) {
		case Integer:
			return left + right, nil
		default:
			return nil, InvalidOperandError{tokens.Plus.String(), right}
		}
	case String:
		switch right := right.(type) {
		case String:
			return left + right, nil
		default:
			return nil, InvalidOperandError{tokens.Plus.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.Plus.String(), left}
	}
}

func evalInfixMinusOperation(left, right Object) (Object, error) {
	leftInteger, ok := left.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), left}
	}

	rightInteger, ok := right.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), right}
	}

	return leftInteger - rightInteger, nil
}

func evalInfixAsteriskOperation(left, right Object) (Object, error) {
	leftInteger, ok := left.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Asterisk.String(), left}
	}

	rightInteger, ok := right.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Asterisk.String(), right}
	}

	return leftInteger * rightInteger, nil
}

func evalInfixSlashOperation(left, right Object) (Object, error) {
	leftInteger, ok := left.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Slash.String(), left}
	}

	rightInteger, ok := right.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Slash.String(), right}
	}

	return leftInteger / rightInteger, nil
}

func evalInfixEqualsOperation(left, right Object) (Object, error) {
	return Boolean(left == right), nil
}

func evalInfixNotEqualsOperation(left, right Object) (Object, error) {
	return Boolean(left != right), nil
}

func evalInfixLessThanOperation(left, right Object) (Object, error) {
	switch left := left.(type) {
	case Integer:
		switch right := right.(type) {
		case Integer:
			return Boolean(left < right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	case String:
		switch right := right.(type) {
		case String:
			return Boolean(left < right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.LessThan.String(), left}
	}
}

func evalInfixGreaterThanOperation(left, right Object) (Object, error) {
	switch left := left.(type) {
	case Integer:
		switch right := right.(type) {
		case Integer:
			return Boolean(left > right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	case String:
		switch right := right.(type) {
		case String:
			return Boolean(left > right), nil
		default:
			return nil, InvalidOperandError{tokens.LessThan.String(), right}
		}
	default:
		return nil, InvalidOperandError{tokens.LessThan.String(), left}
	}
}

func evalInfixLessThanOrEqualsOperation(left, right Object) (Object, error) {
	lessThan, err := evalInfixLessThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return lessThan.(Boolean) || equal.(Boolean), nil
}

func evalInfixGreaterThanOrEqualsOperation(left, right Object) (Object, error) {
	greaterThan, err := evalInfixGreaterThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return greaterThan.(Boolean) || equal.(Boolean), nil
}
