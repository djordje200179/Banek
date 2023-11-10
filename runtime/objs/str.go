package objs

import (
	"banek/runtime/types"
	"strings"
)

type Str string

func (str Str) Type() types.Type { return types.TypeStr }
func (str Str) Clone() types.Obj { return str }
func (str Str) String() string   { return string(str) }

func (str Str) Equals(other types.Obj) bool {
	otherStr, ok := other.(Str)
	if !ok {
		return false
	}

	return str == otherStr
}

func (str Str) Add(other types.Obj) (types.Obj, bool) {
	otherStr, ok := other.(Str)
	if !ok {
		return nil, false
	}

	return str + otherStr, true
}

func (str Str) Mul(other types.Obj) (types.Obj, bool) {
	count, ok := other.(Int)
	if !ok {
		return nil, false
	}

	if count < 0 {
		return Str(""), true
	}

	var sb strings.Builder
	sb.Grow(int(count) * len(str))
	for i := 0; i < int(count); i++ {
		sb.WriteString(string(str))
	}

	return Str(sb.String()), true
}

func (str Str) Less(other types.Obj) (less, ok bool) {
	otherStr, ok := other.(Str)
	if !ok {
		return
	}

	return str < otherStr, true
}
