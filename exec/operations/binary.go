package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"slices"
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
)

func (operator BinaryOperator) String() string {
	return binaryOperatorNames[operator]
}

type binaryOp func(left, right objects.Object) (objects.Object, error)

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

var binaryOps = [...]binaryOp{
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

func EvalBinary(left, right objects.Object, operator BinaryOperator) (objects.Object, error) {
	if operator >= BinaryOperator(len(binaryOps)) {
		return nil, errors.ErrUnknownOperator{Operator: operator.String()}
	}

	return binaryOps[operator](left, right)
}

func evalBinaryPlus(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return left + right, nil
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return left + right, nil
		}
	case objects.Array:
		switch right := right.(type) {
		case objects.Array:
			newArray := make(objects.Array, len(left)+len(right))
			copy(newArray, left)
			copy(newArray[len(left):], right)

			return newArray, nil
		}
	}

	return nil, errors.ErrInvalidOp{Operator: BinaryPlus.String(), LeftOperand: left, RightOperand: right}
}

func evalBinaryMinus(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryMinus.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryMinus.String(), LeftOperand: left, RightOperand: right}
	}

	return leftInteger - rightInteger, nil
}

func evalBinaryAsterisk(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return left * right, nil
		}
	case objects.String:
		switch right := right.(type) {
		case objects.Integer:
			var sb strings.Builder
			sb.Grow(len(left) * int(right))

			for i := 0; i < int(right); i++ {
				sb.WriteString(string(left))
			}

			return objects.String(sb.String()), nil
		}
	case objects.Array:
		switch right := right.(type) {
		case objects.Integer:
			newArray := make(objects.Array, len(left)*int(right))
			for i := 0; i < int(right); i++ {
				copy(newArray[i*len(left):], left)
			}

			return newArray, nil
		}
	}

	return nil, errors.ErrInvalidOp{Operator: BinaryAsterisk.String(), LeftOperand: left, RightOperand: right}
}

func evalBinaryCaret(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryCaret.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryCaret.String(), LeftOperand: left, RightOperand: right}
	}

	result := objects.Integer(1)
	if rightInteger < 0 {
		for i := objects.Integer(0); i > rightInteger; i-- {
			result /= leftInteger
		}
	} else {
		for i := objects.Integer(0); i < rightInteger; i++ {
			result *= leftInteger
		}
	}

	return result, nil
}

func evalBinarySlash(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinarySlash.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinarySlash.String(), LeftOperand: left, RightOperand: right}
	}

	return leftInteger / rightInteger, nil
}

func evalBinaryModulo(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryModulo.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOp{Operator: BinaryModulo.String(), LeftOperand: left, RightOperand: right}
	}

	return leftInteger % rightInteger, nil
}

func evalBinaryEquals(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Array:
		switch right := right.(type) {
		case objects.Array:
			return objects.Boolean(slices.Equal(left, right)), nil
		}
	}

	return objects.Boolean(left == right), nil
}

func evalBinaryNotEquals(left, right objects.Object) (objects.Object, error) {
	res, err := evalBinaryEquals(left, right)
	if err != nil {
		return nil, err
	}

	return !res.(objects.Boolean), nil
}

func evalBinaryLess(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return objects.Boolean(left < right), nil
		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left < right), nil
		}
	}

	return nil, errors.ErrInvalidOp{Operator: BinaryLess.String(), LeftOperand: left, RightOperand: right}
}

func evalBinaryGreater(left, right objects.Object) (objects.Object, error) {
	switch left := left.(type) {
	case objects.Integer:
		switch right := right.(type) {
		case objects.Integer:
			return objects.Boolean(left > right), nil
		default:

		}
	case objects.String:
		switch right := right.(type) {
		case objects.String:
			return objects.Boolean(left > right), nil
		}
	}

	return nil, errors.ErrInvalidOp{Operator: BinaryGreater.String(), LeftOperand: left, RightOperand: right}
}

func evalBinaryLessEquals(left, right objects.Object) (objects.Object, error) {
	less, err := evalBinaryLess(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalBinaryEquals(left, right)
	if err != nil {
		return nil, err
	}

	return less.(objects.Boolean) || equal.(objects.Boolean), nil
}

func evalBinaryGreaterEquals(left, right objects.Object) (objects.Object, error) {
	greater, err := evalBinaryGreater(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalBinaryEquals(left, right)
	if err != nil {
		return nil, err
	}

	return greater.(objects.Boolean) || equal.(objects.Boolean), nil
}
