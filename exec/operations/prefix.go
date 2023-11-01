package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

type PrefixOperationType byte

const (
	PrefixMinusOperation PrefixOperationType = iota
	PrefixBangOperation
)

func (operation PrefixOperationType) String() string {
	return prefixOperationNames[operation]
}

type prefixOperationFunction func(operand objects.Object) (objects.Object, error)

var prefixOperationNames = []string{
	PrefixMinusOperation: "-",
	PrefixBangOperation:  "!",
}

var prefixOperations = []prefixOperationFunction{
	PrefixMinusOperation: evalPrefixMinusOperation,
	PrefixBangOperation:  evalPrefixBangOperation,
}

func EvalPrefixOperation(operand objects.Object, operation PrefixOperationType) (objects.Object, error) {
	if operation >= PrefixOperationType(len(prefixOperationNames)) {
		return nil, errors.ErrUnknownOperator{Operator: operation.String()}
	}

	return prefixOperations[operation](operand)
}

func evalPrefixMinusOperation(operand objects.Object) (objects.Object, error) {
	integer, ok := operand.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: PrefixMinusOperation.String(), RightOperand: operand}
	}

	return -integer, nil
}

func evalPrefixBangOperation(operand objects.Object) (objects.Object, error) {
	boolean, ok := operand.(objects.Boolean)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: PrefixBangOperation.String(), RightOperand: operand}
	}

	return !boolean, nil
}
