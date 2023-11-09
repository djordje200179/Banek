package ops

import (
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/types"
)

type BinaryOperator byte

const (
	BinaryPlus BinaryOperator = iota
	BinaryMinus
	BinaryAsterisk
	BinarySlash
	BinaryModulo
	BinaryCaret

	BinaryEquals
	BinaryNotEquals
	BinaryLess
	BinaryGreater
	BinaryLessEquals
	BinaryGreaterEquals
)

func (operator BinaryOperator) String() string {
	return binaryOperatorNames[operator]
}

type binaryOp func(left, right types.Obj) (types.Obj, error)

var binaryOperatorNames = [...]string{
	BinaryPlus:     "+",
	BinaryMinus:    "-",
	BinaryAsterisk: "*",
	BinarySlash:    "/",
	BinaryModulo:   "%",
	BinaryCaret:    "^",

	BinaryEquals:        "==",
	BinaryNotEquals:     "!=",
	BinaryLess:          "<",
	BinaryGreater:       ">",
	BinaryLessEquals:    "<=",
	BinaryGreaterEquals: ">=",
}

var BinaryOps = [...]binaryOp{
	BinaryPlus:     evalBinaryPlus,
	BinaryMinus:    evalBinaryMinus,
	BinaryAsterisk: evalBinaryAsterisk,
	BinarySlash:    evalBinarySlash,
	BinaryModulo:   evalBinaryModulo,
	BinaryCaret:    evalBinaryCaret,

	BinaryEquals:        evalBinaryEquals,
	BinaryNotEquals:     evalBinaryNotEquals,
	BinaryLess:          evalBinaryLess,
	BinaryGreater:       evalBinaryGreater,
	BinaryLessEquals:    evalBinaryLessEquals,
	BinaryGreaterEquals: evalBinaryGreaterEquals,
}

func evalBinaryPlus(left, right types.Obj) (types.Obj, error) {
	leftAdder, ok := left.(types.Adder)
	if !ok || !leftAdder.CanAdd(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryPlus.String(), LeftOperand: left, RightOperand: right}
	}

	return leftAdder.Add(right), nil
}

func evalBinaryMinus(left, right types.Obj) (types.Obj, error) {
	leftAdder, ok := left.(types.Subber)
	if !ok || !leftAdder.CanSub(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryMinus.String(), LeftOperand: left, RightOperand: right}
	}

	return leftAdder.Sub(right), nil
}

func evalBinaryAsterisk(left, right types.Obj) (types.Obj, error) {
	leftAdder, ok := left.(types.Multer)
	if !ok || !leftAdder.CanMultiply(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryAsterisk.String(), LeftOperand: left, RightOperand: right}
	}

	return leftAdder.Multiply(right), nil
}

func evalBinaryCaret(left, right types.Obj) (types.Obj, error) {
	leftPowwer, ok := left.(types.Powwer)
	if !ok || !leftPowwer.CanPow(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryCaret.String(), LeftOperand: left, RightOperand: right}
	}

	return leftPowwer.Pow(right), nil
}

func evalBinarySlash(left, right types.Obj) (types.Obj, error) {
	leftAdder, ok := left.(types.Diver)
	if !ok || !leftAdder.CanDivide(right) {
		return nil, errors.ErrInvalidOp{Operator: BinarySlash.String(), LeftOperand: left, RightOperand: right}
	}

	return leftAdder.Divide(right), nil
}

func evalBinaryModulo(left, right types.Obj) (types.Obj, error) {
	leftModder, ok := left.(types.Modder)
	if !ok || !leftModder.CanMod(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryModulo.String(), LeftOperand: left, RightOperand: right}
	}

	return leftModder.Mod(right), nil
}

func evalBinaryEquals(left, right types.Obj) (types.Obj, error) {
	equality := left.Equals(right)

	return objs.Bool(equality), nil
}

func evalBinaryNotEquals(left, right types.Obj) (types.Obj, error) {
	equality := left.Equals(right)

	return objs.Bool(!equality), nil
}

func evalBinaryLess(left, right types.Obj) (types.Obj, error) {
	leftLesser, ok := left.(types.Lesser)
	if ok && !leftLesser.CanLess(right) {
		return nil, errors.ErrInvalidOp{Operator: BinaryLess.String(), LeftOperand: left, RightOperand: right}
	}

	isLess := leftLesser.Less(right)

	return objs.Bool(isLess), nil
}

func evalBinaryGreater(left, right types.Obj) (types.Obj, error) {
	return evalBinaryLess(right, left)
}

func evalBinaryLessEquals(left, right types.Obj) (types.Obj, error) {
	less, err := evalBinaryLess(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalBinaryEquals(left, right)
	if err != nil {
		return nil, err
	}

	return less.(objs.Bool) || equal.(objs.Bool), nil
}

func evalBinaryGreaterEquals(left, right types.Obj) (types.Obj, error) {
	return evalBinaryLessEquals(right, left)
}
