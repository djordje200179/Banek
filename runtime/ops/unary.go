package ops

import (
	"banek/runtime/objs"
	"strings"
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

type unaryOp func(operand objs.Obj) (objs.Obj, error)

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

type ErrInvalidUnaryOpOperand struct {
	Operator UnaryOperator

	Operand objs.Obj
}

func (err ErrInvalidUnaryOpOperand) Error() string {
	var sb strings.Builder

	sb.WriteString("invalid operand for ")
	sb.WriteString(err.Operator.String())
	sb.WriteString(": ")
	sb.WriteString(err.Operand.String())

	return sb.String()
}

func evalUnaryMinus(operand objs.Obj) (objs.Obj, error) {
	switch operand.Tag {
	case objs.TypeInt:
		integer := operand.AsInt()
		return objs.MakeInt(-integer), nil
	default:
		return objs.Obj{}, ErrInvalidUnaryOpOperand{Operator: UnaryMinus, Operand: operand}
	}
}

func evalUnaryBang(operand objs.Obj) (objs.Obj, error) {
	switch operand.Tag {
	case objs.TypeBool:
		boolean := operand.AsBool()
		return objs.MakeBool(!boolean), nil
	default:
		return objs.Obj{}, ErrInvalidUnaryOpOperand{Operator: UnaryBang, Operand: operand}
	}
}

func evalUnaryLeftArrow(operand objs.Obj) (objs.Obj, error) {
	switch operand.Tag {
	case objs.TypeArray:
		arr := operand.AsArray()

		if len(arr.Slice) == 0 {
			return objs.Obj{}, nil
		}

		elem := arr.Slice[0]
		arr.Slice = arr.Slice[1:]

		return elem, nil
	default:
		return objs.Obj{}, ErrInvalidUnaryOpOperand{Operator: UnaryLeftArrow, Operand: operand}
	}
}
