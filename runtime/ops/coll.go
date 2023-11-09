package ops

import (
	"banek/runtime/errors"
	"banek/runtime/types"
)

func EvalCollSet(coll, key, value types.Obj) error {
	collColl, ok := coll.(types.Coll)
	if !ok || !collColl.CanIndex(key) {
		return errors.ErrInvalidOp{Operator: "[]", LeftOperand: coll, RightOperand: key}
	}

	return collColl.Set(key, value)
}

func EvalCollGet(coll, key types.Obj) (types.Obj, error) {
	collColl, ok := coll.(types.Coll)
	if !ok || !collColl.CanIndex(key) {
		return nil, errors.ErrInvalidOp{Operator: "[]", LeftOperand: coll, RightOperand: key}
	}

	return collColl.Get(key)
}
