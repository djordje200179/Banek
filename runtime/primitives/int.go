package primitives

import (
	"banek/runtime"
	"cmp"
	"fmt"
)

type Int int

func (i Int) String() string                { return fmt.Sprintf("%d", i) }
func (i Int) Truthy() bool                  { return i != 0 }
func (i Int) Clone() runtime.Obj            { return i }
func (i Int) Equals(other runtime.Obj) bool { return i == other }

func (i Int) Add(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return nil, false
	}

	return i + otherInt, true
}

func (i Int) Sub(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return nil, false
	}

	return i - otherInt, true
}

func (i Int) Mul(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return nil, false
	}

	return i * otherInt, true
}

func (i Int) Div(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return nil, false
	}

	return i / otherInt, true
}

func (i Int) Mod(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return nil, false
	}

	return i % otherInt, true
}

func (i Int) Compare(other runtime.Obj) (int, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok {
		return 0, false
	}

	return cmp.Compare(i, otherInt), true
}

func (i Int) Neg() (runtime.Obj, bool) {
	return -i, true
}
