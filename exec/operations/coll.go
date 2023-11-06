package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

func EvalCollSet(coll, key, value objects.Object) error {
	switch coll := coll.(type) {
	case objects.Array:
		index, ok := key.(objects.Integer)
		if !ok {
			return errors.ErrInvalidOp{Operator: "Index", LeftOperand: coll, RightOperand: key}
		}

		if index < 0 {
			index += objects.Integer(len(coll))
		}

		if index < 0 || index >= objects.Integer(len(coll)) {
			return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(coll)}
		}

		coll[index] = value

		return nil
	default:
		return errors.ErrInvalidOp{Operator: "index", LeftOperand: coll, RightOperand: key}
	}
}

func EvalCollGet(coll, key objects.Object) (objects.Object, error) {
	switch coll := coll.(type) {
	case objects.Array:
		index, ok := key.(objects.Integer)
		if !ok {
			return nil, errors.ErrInvalidOp{Operator: "index", LeftOperand: coll, RightOperand: key}
		}

		if index < 0 {
			index += objects.Integer(len(coll))
		}

		if index < 0 || index >= objects.Integer(len(coll)) {
			return nil, objects.ErrIndexOutOfBounds{Index: int(index), Size: len(coll)}
		}

		return coll[index], nil
	default:
		return nil, errors.ErrInvalidOp{Operator: "Index", LeftOperand: coll, RightOperand: key}
	}
}
