package ops

import (
	"banek/runtime/errors"
	"banek/runtime/types"
)

type UnaryOperator byte

const (
	UnaryMinus UnaryOperator = iota
	UnaryBang
)

func (operation UnaryOperator) String() string {
	return unaryOperatorNames[operation]
}

type unaryOp func(operand types.Obj) (types.Obj, error)

var unaryOperatorNames = [...]string{
	UnaryMinus: "-",
	UnaryBang:  "!",
}

var UnaryOps = [...]unaryOp{
	UnaryMinus: evalUnaryMinus,
	UnaryBang:  evalUnaryBang,
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
