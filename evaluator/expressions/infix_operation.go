package expressions

import (
	"banek/evaluator/objects"
	"banek/tokens"
)

type infixOperationEvaluator func(objects.Object, objects.Object) (objects.Object, error)

var infixOperations = map[tokens.TokenType]infixOperationEvaluator{
	tokens.Plus:     evalInfixPlusOperation,
	tokens.Minus:    evalInfixMinusOperation,
	tokens.Asterisk: evalInfixAsteriskOperation,
	tokens.Slash:    evalInfixSlashOperation,

	tokens.Equals:              evalInfixEqualsOperation,
	tokens.NotEquals:           evalInfixNotEqualsOperation,
	tokens.LessThan:            evalInfixLessThanOperation,
	tokens.GreaterThan:         evalInfixGreaterThanOperation,
	tokens.LessThanOrEquals:    evalInfixLessThanOrEqualsOperation,
	tokens.GreaterThanOrEquals: evalInfixGreaterThanOrEqualsOperation,
}

func evalInfixOperation(operator tokens.Token, left, right objects.Object) (objects.Object, error) {
	operation := infixOperations[operator.Type]
	if operation == nil {
		return nil, UnknownOperatorError{operator.Type}
	}

	return operation(left, right)
}

func evalInfixPlusOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Plus.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Plus.String(), right}
	}

	return leftInteger + rightInteger, nil
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
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.LessThan.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.LessThan.String(), right}
	}

	return objects.Boolean(leftInteger < rightInteger), nil
}

func evalInfixGreaterThanOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.GreaterThan.String(), left}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.GreaterThan.String(), right}
	}

	return objects.Boolean(leftInteger > rightInteger), nil
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
