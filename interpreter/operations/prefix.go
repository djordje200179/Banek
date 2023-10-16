package operations

import (
	errors2 "banek/exec/errors"
	objects2 "banek/exec/objects"
	"banek/tokens"
)

type prefixOperation func(operand objects2.Object) (objects2.Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperation{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixBangOperation,
}

func EvalPrefixOperation(operand objects2.Object, operator tokens.TokenType) (objects2.Object, error) {
	operation := prefixOperations[operator]
	if operation == nil {
		return nil, errors2.UnknownOperatorError{Operator: operator}
	}

	return operation(operand)
}

func evalPrefixMinusOperation(operand objects2.Object) (objects2.Object, error) {
	integer, ok := operand.(objects2.Integer)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Minus.String(), Operand: operand}
	}

	return -integer, nil
}

func evalPrefixBangOperation(operand objects2.Object) (objects2.Object, error) {
	boolean, ok := operand.(objects2.Boolean)
	if !ok {
		return nil, errors2.InvalidOperandError{Operator: tokens.Bang.String(), Operand: operand}
	}

	return !boolean, nil
}
