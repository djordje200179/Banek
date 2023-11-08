package operations

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

func EvalCollSet(coll, key, value objects.Object) error {
	collColl, ok := coll.(objects.Coll)
	if !ok || !collColl.AcceptsKey(key) {
		return errors.ErrInvalidOp{Operator: "[]", LeftOperand: coll, RightOperand: key}
	}

	return collColl.Set(key, value)
}

func EvalCollGet(coll, key objects.Object) (objects.Object, error) {
	collColl, ok := coll.(objects.Coll)
	if !ok || !collColl.AcceptsKey(key) {
		return nil, errors.ErrInvalidOp{Operator: "[]", LeftOperand: coll, RightOperand: key}
	}

	return collColl.Get(key)
}
