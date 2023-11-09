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

func (str Str) CanAdd(other types.Obj) bool {
	_, ok := other.(Str)
	return ok
}

func (str Str) Add(other types.Obj) types.Obj {
	otherStr := other.(Str)
	return str + otherStr
}

func (str Str) CanMultiply(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (str Str) Multiply(other types.Obj) types.Obj {
	count := other.(Int)
	if count < 0 {
		return Str("")
	}

	var sb strings.Builder
	for i := 0; i < int(count); i++ {
		sb.WriteString(string(str))
	}

	return Str(sb.String())
}

func (str Str) CanLess(other types.Obj) bool {
	_, ok := other.(Str)
	return ok
}

func (str Str) Less(other types.Obj) bool {
	otherStr := other.(Str)
	return str < otherStr
}
