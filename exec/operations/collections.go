package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

func EvalCollectionSet(collection, key, value objects.Object) error {
	switch collection := collection.(type) {
	case objects.Array:
		index, ok := key.(objects.Integer)
		if !ok {
			return errors.ErrInvalidOperand{Operation: "Index", LeftOperand: collection, RightOperand: key}
		}

		if index < 0 {
			index += objects.Integer(len(collection))
		}

		if index < 0 || index >= objects.Integer(len(collection)) {
			return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(collection)}
		}

		collection[index] = value

		return nil
	default:
		return errors.ErrInvalidOperand{Operation: "Index", LeftOperand: collection, RightOperand: key}
	}
}

func EvalCollectionGet(collection, key objects.Object) (objects.Object, error) {
	switch collection := collection.(type) {
	case objects.Array:
		index, ok := key.(objects.Integer)
		if !ok {
			return nil, errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: key}
		}

		if index < 0 {
			index += objects.Integer(len(collection))
		}

		if index < 0 || index >= objects.Integer(len(collection)) {
			return nil, objects.ErrIndexOutOfBounds{Index: int(index), Size: len(collection)}
		}

		return collection[index], nil
	default:
		return nil, errors.ErrInvalidOperand{Operation: "Index", LeftOperand: collection, RightOperand: key}
	}
}
