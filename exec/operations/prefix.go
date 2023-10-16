package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/tokens"
)

type prefixOperation func(operand objects.Object) (objects.Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperation{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixBangOperation,
}

func EvalPrefixOperation(operand objects.Object, operator tokens.TokenType) (objects.Object, error) {
	operation := prefixOperations[operator]
	if operation == nil {
		return nil, errors.UnknownOperatorError{Operator: operator}
	}

	return operation(operand)
}

func evalPrefixMinusOperation(operand objects.Object) (objects.Object, error) {
	integer, ok := operand.(objects.Integer)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Minus.String(), Operand: operand}
	}

	return -integer, nil
}

func evalPrefixBangOperation(operand objects.Object) (objects.Object, error) {
	boolean, ok := operand.(objects.Boolean)
	if !ok {
		return nil, errors.InvalidOperandError{Operator: tokens.Bang.String(), Operand: operand}
	}

	return !boolean, nil
}
