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
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryPlus.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftAdder.Add(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryPlus.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
}

func evalBinaryMinus(left, right types.Obj) (types.Obj, error) {
	leftSubber, ok := left.(types.Subber)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryMinus.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftSubber.Sub(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryMinus.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
}

func evalBinaryAsterisk(left, right types.Obj) (types.Obj, error) {
	leftMulter, ok := left.(types.Multer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryAsterisk.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftMulter.Mul(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryAsterisk.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
}

func evalBinaryCaret(left, right types.Obj) (types.Obj, error) {
	leftPowwer, ok := left.(types.Powwer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryCaret.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftPowwer.Pow(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryCaret.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
}

func evalBinarySlash(left, right types.Obj) (types.Obj, error) {
	leftDiver, ok := left.(types.Diver)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinarySlash.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftDiver.Div(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinarySlash.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
}

func evalBinaryModulo(left, right types.Obj) (types.Obj, error) {
	leftModder, ok := left.(types.Modder)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryModulo.String(), LeftOperand: left, RightOperand: right}
	}

	res, ok := leftModder.Mod(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryModulo.String(), LeftOperand: left, RightOperand: right}
	}

	return res, nil
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
	if ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryLess.String(), LeftOperand: left, RightOperand: right}
	}

	isLess, ok := leftLesser.Less(right)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryLess.String(), LeftOperand: left, RightOperand: right}
	}

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
