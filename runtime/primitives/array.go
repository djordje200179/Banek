package primitives

import (
	"banek/runtime"
	"slices"
	"strings"
)

type Array []runtime.Obj

func (arr Array) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, obj := range arr {
		sb.WriteString(obj.String())
		if i != len(arr)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")

	return sb.String()
}

func (arr Array) Truthy() bool       { return len(arr) > 0 }
func (arr Array) Clone() runtime.Obj { return slices.Clone(arr) }

func (arr Array) Equals(other runtime.Obj) bool {
	var otherArr Array
	var ok bool
	if otherArr, ok = other.(Array); !ok {
		return false
	}

	return slices.Equal(arr, otherArr)
}

func (arr Array) Add(other runtime.Obj) (runtime.Obj, bool) {
	var otherArr Array
	var ok bool
	if otherArr, ok = other.(Array); !ok {
		return nil, false
	}

	res := make(Array, len(arr)+len(otherArr))
	copy(res[:len(arr)], arr)
	copy(res[len(arr):], otherArr)

	return res, true
}

func (arr Array) Mul(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok || otherInt < 0 {
		return nil, false
	}

	res := make(Array, len(arr)*int(otherInt))
	for i := range int(otherInt) {
		copy(res[i*len(arr):(i+1)*len(arr)], arr)
	}

	return res, true
}
