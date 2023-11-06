package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

type UnaryOperator byte

const (
	UnaryMinus UnaryOperator = iota
	UnaryBang
)

func (operation UnaryOperator) String() string {
	return unaryOperatorNames[operation]
}

type unaryOp func(operand objects.Object) (objects.Object, error)

var unaryOperatorNames = [...]string{
	UnaryMinus: "-",
	UnaryBang:  "!",
}

var unaryOps = [...]unaryOp{
	UnaryMinus: evalUnaryMinus,
	UnaryBang:  evalUnaryBang,
}

func EvalUnary(operand objects.Object, operator UnaryOperator) (objects.Object, error) {
	if operator >= UnaryOperator(len(unaryOps)) {
		return nil, errors.ErrUnknownOperator{Operator: operator.String()}
	}

	return unaryOps[operator](operand)
}

func evalUnaryMinus(operand objects.Object) (objects.Object, error) {
	integer, ok := operand.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: UnaryMinus.String(), RightOperand: operand}
	}

	return -integer, nil
}

func evalUnaryBang(operand objects.Object) (objects.Object, error) {
	boolean, ok := operand.(objects.Boolean)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: UnaryBang.String(), RightOperand: operand}
	}

	return !boolean, nil
}
