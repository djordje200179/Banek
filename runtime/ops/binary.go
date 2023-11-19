package ops

import (
	"banek/runtime/objs"
	"strings"
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

	BinaryLeftArrow
)

func (operator BinaryOperator) String() string {
	return binaryOperatorNames[operator]
}

type binaryOp func(left, right objs.Obj) (objs.Obj, error)

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

	BinaryLeftArrow: "<-",
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

	BinaryLeftArrow: evalBinaryLeftArrow,
}

type ErrInvalidBinaryOpOperand struct {
	Operator BinaryOperator

	LeftOperand, RightOperand objs.Obj
}

func (err ErrInvalidBinaryOpOperand) Error() string {
	var sb strings.Builder

	sb.WriteString("invalid operands for ")
	sb.WriteString(err.Operator.String())
	sb.WriteString(": ")
	sb.WriteString(err.LeftOperand.String())
	sb.WriteString(" and ")
	sb.WriteString(err.RightOperand.String())

	return sb.String()
}

func evalBinaryPlus(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != right.Tag {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryPlus, LeftOperand: left, RightOperand: right}
	}

	switch left.Tag {
	case objs.TypeInt:
		return objs.MakeInt(left.AsInt() + right.AsInt()), nil
	case objs.TypeStr:
		return objs.MakeStr(left.AsStr() + right.AsStr()), nil
	case objs.TypeArray:
		leftArr := left.AsArray()
		rightArr := right.AsArray()

		newArr := new(objs.Array)
		newArr.Slice = make([]objs.Obj, len(leftArr.Slice)+len(rightArr.Slice))
		copy(newArr.Slice, leftArr.Slice)
		copy(newArr.Slice[len(leftArr.Slice):], rightArr.Slice)

		return objs.MakeArray(newArr), nil
	default:
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryPlus, LeftOperand: left, RightOperand: right}
	}
}

func evalBinaryMinus(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != objs.TypeInt || right.Tag != objs.TypeInt {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryMinus, LeftOperand: left, RightOperand: right}
	}

	return objs.MakeInt(left.AsInt() - right.AsInt()), nil
}

func evalBinaryAsterisk(left, right objs.Obj) (objs.Obj, error) {
	switch left.Tag {
	case objs.TypeInt:
		if right.Tag != objs.TypeInt {
			return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryAsterisk, LeftOperand: left, RightOperand: right}
		}

		return objs.MakeInt(left.AsInt() * right.AsInt()), nil
	case objs.TypeArray:
		if right.Tag != objs.TypeInt {
			return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryAsterisk, LeftOperand: left, RightOperand: right}
		}

		baseArr := left.AsArray()

		newArr := new(objs.Array)
		newArr.Slice = make([]objs.Obj, len(baseArr.Slice)*right.AsInt())

		for i := 0; i < right.AsInt(); i++ {
			copy(newArr.Slice[i*len(baseArr.Slice):], baseArr.Slice)
		}

		return objs.MakeArray(newArr), nil
	case objs.TypeStr:
		if right.Tag != objs.TypeInt {
			return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryAsterisk, LeftOperand: left, RightOperand: right}
		}

		baseStr := left.AsStr()

		var sb strings.Builder
		sb.Grow(len(baseStr) * right.AsInt())

		for i := 0; i < right.AsInt(); i++ {
			sb.WriteString(baseStr)
		}

		return objs.MakeStr(sb.String()), nil
	default:
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryAsterisk, LeftOperand: left, RightOperand: right}
	}
}

func evalBinaryCaret(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != objs.TypeInt || right.Tag != objs.TypeInt {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryCaret, LeftOperand: left, RightOperand: right}
	}

	base := left.AsInt()
	power := right.AsInt()

	if power < 0 {
		if base == 1 {
			return objs.MakeInt(1), nil
		}

		return objs.MakeInt(0), nil
	}

	result := 1
	for i := 0; i < power; i++ {
		result *= base
	}

	return objs.MakeInt(result), nil
}

func evalBinarySlash(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != objs.TypeInt || right.Tag != objs.TypeInt {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinarySlash, LeftOperand: left, RightOperand: right}
	}

	return objs.MakeInt(left.AsInt() / right.AsInt()), nil
}

func evalBinaryModulo(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != objs.TypeInt || right.Tag != objs.TypeInt {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryModulo, LeftOperand: left, RightOperand: right}
	}

	return objs.MakeInt(left.AsInt() % right.AsInt()), nil
}

func evalBinaryEquals(left, right objs.Obj) (objs.Obj, error) {
	return objs.MakeBool(left.Equals(right)), nil
}

func evalBinaryNotEquals(left, right objs.Obj) (objs.Obj, error) {
	return objs.MakeBool(!left.Equals(right)), nil
}

func evalBinaryLess(left, right objs.Obj) (objs.Obj, error) {
	if left.Tag != right.Tag {
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryLess, LeftOperand: left, RightOperand: right}
	}

	switch left.Tag {
	case objs.TypeInt:
		return objs.MakeBool(left.AsInt() < right.AsInt()), nil
	case objs.TypeStr:
		return objs.MakeBool(left.AsStr() < right.AsStr()), nil
	default:
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryLess, LeftOperand: left, RightOperand: right}
	}
}

func evalBinaryGreater(left, right objs.Obj) (objs.Obj, error) {
	return evalBinaryLess(right, left)
}

func evalBinaryLessEquals(left, right objs.Obj) (objs.Obj, error) {
	equal, _ := evalBinaryEquals(left, right)
	if equal.AsBool() {
		return equal, nil
	}

	return evalBinaryLess(left, right)
}

func evalBinaryGreaterEquals(left, right objs.Obj) (objs.Obj, error) {
	return evalBinaryLessEquals(right, left)
}

func evalBinaryLeftArrow(left, right objs.Obj) (objs.Obj, error) {
	switch left.Tag {
	case objs.TypeArray:
		arr := left.AsArray()

		arr.Slice = append(arr.Slice, right)

		return objs.MakeArray(arr), nil
	default:
		return objs.Obj{}, ErrInvalidBinaryOpOperand{Operator: BinaryLeftArrow, LeftOperand: left, RightOperand: right}
	}
}
