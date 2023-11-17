package ops

import (
	"banek/runtime/objs"
	"fmt"
)

type ErrNotIndexable struct {
	Coll, Key objs.Obj
}

func (err ErrNotIndexable) Error() string {
	return fmt.Sprintf("not indexable: %s[%s]", err.Coll, err.Key)
}

type ErrIndexOutOfBounds struct {
	Index int
	Size  int
}

func (err ErrIndexOutOfBounds) Error() string {
	return fmt.Sprintf("index out of bounds: index %d, size %d", err.Index, err.Size)
}

func EvalCollSet(coll, key, value objs.Obj) error {
	switch coll.Tag {
	case objs.TypeArray:
		arr := coll.AsArray()

		if key.Tag != objs.TypeInt {
			return ErrNotIndexable{Coll: coll, Key: key}
		}

		index := key.AsInt()
		if index < 0 {
			index += len(arr.Slice)
		}

		if index < 0 || index >= len(arr.Slice) {
			return ErrIndexOutOfBounds{Index: index, Size: len(arr.Slice)}
		}

		arr.Slice[index] = value

		return nil
	default:
		return ErrNotIndexable{Coll: coll, Key: key}
	}
}

func EvalCollGet(coll, key objs.Obj) (objs.Obj, error) {
	switch coll.Tag {
	case objs.TypeArray:
		arr := coll.AsArray()

		if key.Tag != objs.TypeInt {
			return objs.Obj{}, ErrNotIndexable{Coll: coll, Key: key}
		}

		index := key.AsInt()
		if index < 0 {
			index += len(arr.Slice)
		}

		if index < 0 || index >= len(arr.Slice) {
			return objs.Obj{}, ErrIndexOutOfBounds{Index: index, Size: len(arr.Slice)}
		}

		return arr.Slice[index], nil
	default:
		return objs.Obj{}, ErrNotIndexable{Coll: coll, Key: key}
	}
}
