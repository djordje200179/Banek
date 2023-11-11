package ops

import (
	"banek/runtime/errors"
	"banek/runtime/types"
)

type UnaryOperator byte

const (
	UnaryMinus UnaryOperator = iota
	UnaryBang
	UnaryLeftArrow
)

func (operation UnaryOperator) String() string {
	return unaryOperatorNames[operation]
}

type unaryOp func(operand types.Obj) (types.Obj, error)

var unaryOperatorNames = [...]string{
	UnaryMinus:     "-",
	UnaryBang:      "!",
	UnaryLeftArrow: "<-",
}

var UnaryOps = [...]unaryOp{
	UnaryMinus:     evalUnaryMinus,
	UnaryBang:      evalUnaryBang,
	UnaryLeftArrow: evalUnaryLeftArrow,
}

func evalUnaryMinus(operand types.Obj) (types.Obj, error) {
	negater, ok := operand.(types.Negater)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: UnaryMinus.String(), RightOperand: operand}
	}

	return negater.Negate(), nil
}

func evalUnaryBang(operand types.Obj) (types.Obj, error) {
	notter, ok := operand.(types.Notter)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: UnaryBang.String(), RightOperand: operand}
	}

	return notter.Not(), nil
}

func evalUnaryLeftArrow(operand types.Obj) (types.Obj, error) {
	giver, ok := operand.(types.Giver)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: UnaryLeftArrow.String(), RightOperand: operand}
	}

	return giver.Give(), nil
}
