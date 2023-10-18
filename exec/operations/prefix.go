package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/tokens"
)

type prefixOperation func(operand objects.Object) (objects.Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperation{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixMinusOperation,
}

func EvalPrefixOperation(operand objects.Object, operator tokens.TokenType) (objects.Object, error) {
	operation := prefixOperations[operator]
	if operation == nil {
		return nil, errors.ErrUnknownOperator{Operator: operator}
	}

	return operation(operand)
}

func evalPrefixMinusOperation(operand objects.Object) (objects.Object, error) {
	switch operand := operand.(type) {
	case objects.Integer:
		return -operand, nil
	case objects.Boolean:
		return !operand, nil
	default:
		return nil, errors.ErrInvalidOperand{Operator: tokens.Minus.String(), Operand: operand}
	}
}
