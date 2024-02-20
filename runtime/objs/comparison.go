package objs

import (
	"slices"
	"strings"
)

func (o Obj) Equals(other Obj) bool {
	if o.Type != other.Type {
		return false
	}

	switch o.Type {
	case Int:
		return o.Int == other.Int
	case Bool:
		return o.Int == other.Int
	case String:
		return o.AsString() == other.AsString()
	case Array:
		return slices.EqualFunc(o.AsArray(), other.AsArray(), Obj.Equals)
	case Func:
		return o.Ptr == other.Ptr
	case Builtin:
		return o.Ptr == other.Ptr
	default:
		return true
	}
}

type NotComparableError struct {
	Left, Right Obj
}

func (e NotComparableError) Error() string {
	return "not comparable: " + e.Left.String() + " and " + e.Right.String()
}

func (o Obj) Compare(other Obj) (int, error) {
	err := NotComparableError{o, other}

	if o.Type != other.Type {
		return 0, err
	}

	switch o.Type {
	case Int:
		return o.Int - other.Int, nil
	case String:
		return strings.Compare(o.AsString(), other.AsString()), nil
	default:
		return 0, err
	}
}
