package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"strings"
)

type InfixOperationType byte

const (
	InfixPlusOperation InfixOperationType = iota
	InfixMinusOperation
	InfixAsteriskOperation
	InfixSlashOperation

	InfixEqualsOperation
	InfixNotEqualsOperation
	InfixLessThanOperation
	InfixGreaterThanOperation
	InfixLessThanOrEqualsOperation
	InfixGreaterThanOrEqualsOperation
)

func (operation InfixOperationType) String() string {
	return infixOperationNames[operation]
}

type infixOperationFunction func(left, right objects.Object) (objects.Object, error)

var infixOperationNames = []string{
	InfixPlusOperation:     "+",
	InfixMinusOperation:    "-",
	InfixAsteriskOperation: "*",
	InfixSlashOperation:    "/",

	InfixEqualsOperation:              "==",
	InfixNotEqualsOperation:           "!=",
	InfixLessThanOperation:            "<",
	InfixGreaterThanOperation:         ">",
	InfixLessThanOrEqualsOperation:    "<=",
	InfixGreaterThanOrEqualsOperation: ">=",
}

var infixOperations = []infixOperationFunction{
	InfixPlusOperation:     evalInfixPlusOperation,
	InfixMinusOperation:    evalInfixMinusOperation,
	InfixAsteriskOperation: evalInfixAsteriskOperation,
	InfixSlashOperation:    evalInfixSlashOperation,

	InfixEqualsOperation:              evalInfixEqualsOperation,
	InfixNotEqualsOperation:           evalInfixNotEqualsOperation,
	InfixLessThanOperation:            evalInfixLessThanOperation,
	InfixGreaterThanOperation:         evalInfixGreaterThanOperation,
	InfixLessThanOrEqualsOperation:    evalInfixLessThanOrEqualsOperation,
	InfixGreaterThanOrEqualsOperation: evalInfixGreaterThanOrEqualsOperation,
}

func EvalInfixOperation(left, right objects.Object, operation InfixOperationType) (objects.Object, error) {
	if operation >= InfixOperationType(len(infixOperations)) {
		return nil, errors.ErrUnknownOperator{Operator: operation.String()}
	}

	return infixOperations[operation](left, right)
}

func evalInfixPlusOperation(left, right objects.Object) (objects.Object, error) {
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

	return nil, errors.ErrInvalidOperand{Operation: InfixPlusOperation.String(), LeftOperand: left, RightOperand: right}
}

func evalInfixMinusOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: InfixMinusOperation.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: InfixMinusOperation.String(), LeftOperand: left, RightOperand: right}
	}

	return leftInteger - rightInteger, nil
}

func evalInfixAsteriskOperation(left, right objects.Object) (objects.Object, error) {
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

	return nil, errors.ErrInvalidOperand{Operation: InfixAsteriskOperation.String(), LeftOperand: left, RightOperand: right}
}

func evalInfixSlashOperation(left, right objects.Object) (objects.Object, error) {
	leftInteger, ok := left.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: InfixSlashOperation.String(), LeftOperand: left, RightOperand: right}
	}

	rightInteger, ok := right.(objects.Integer)
	if !ok {
		return nil, errors.ErrInvalidOperand{Operation: InfixSlashOperation.String(), LeftOperand: left, RightOperand: right}
	}

	return leftInteger / rightInteger, nil
}

func evalInfixEqualsOperation(left, right objects.Object) (objects.Object, error) {
	return objects.Boolean(left == right), nil
}

func evalInfixNotEqualsOperation(left, right objects.Object) (objects.Object, error) {
	return objects.Boolean(left != right), nil
}

func evalInfixLessThanOperation(left, right objects.Object) (objects.Object, error) {
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

	return nil, errors.ErrInvalidOperand{Operation: InfixLessThanOperation.String(), LeftOperand: left, RightOperand: right}
}

func evalInfixGreaterThanOperation(left, right objects.Object) (objects.Object, error) {
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

	return nil, errors.ErrInvalidOperand{Operation: InfixGreaterThanOperation.String(), LeftOperand: left, RightOperand: right}
}

func evalInfixLessThanOrEqualsOperation(left, right objects.Object) (objects.Object, error) {
	lessThan, err := evalInfixLessThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return lessThan.(objects.Boolean) || equal.(objects.Boolean), nil
}

func evalInfixGreaterThanOrEqualsOperation(left, right objects.Object) (objects.Object, error) {
	greaterThan, err := evalInfixGreaterThanOperation(left, right)
	if err != nil {
		return nil, err
	}

	equal, err := evalInfixEqualsOperation(left, right)
	if err != nil {
		return nil, err
	}

	return greaterThan.(objects.Boolean) || equal.(objects.Boolean), nil
}
