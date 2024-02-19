package primitives

import (
	"banek/runtime"
	"cmp"
	"strings"
)

type String string

func (s String) String() string                { return string(s) }
func (s String) Truthy() bool                  { return len(s) > 0 }
func (s String) Clone() runtime.Obj            { return s }
func (s String) Equals(other runtime.Obj) bool { return s == other }

func (s String) Add(other runtime.Obj) (runtime.Obj, bool) {
	var otherStr String
	var ok bool
	if otherStr, ok = other.(String); !ok {
		return nil, false
	}

	return s + otherStr, true
}

func (s String) Mul(other runtime.Obj) (runtime.Obj, bool) {
	var otherInt Int
	var ok bool
	if otherInt, ok = other.(Int); !ok || otherInt < 0 {
		return nil, false
	}

	var sb strings.Builder
	for i := 0; i < int(otherInt); i++ {
		sb.WriteString(string(s))
	}

	return String(sb.String()), true
}

func (s String) Compare(other runtime.Obj) (int, bool) {
	var otherStr String
	var ok bool
	if otherStr, ok = other.(String); !ok {
		return 0, false
	}

	return cmp.Compare(s, otherStr), true
}
