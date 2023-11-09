package bytecode

import (
	"banek/runtime/types"
	"slices"
	"strconv"
)

type Func struct {
	TemplateIndex int

	Captures []*types.Obj
}

func (function *Func) Type() types.Type { return types.TypeFunc }
func (function *Func) Clone() types.Obj { return function }
func (function *Func) String() string   { return "func#" + strconv.Itoa(function.TemplateIndex) }

func (function *Func) Equals(other types.Obj) bool {
	otherFunc, ok := other.(*Func)
	if !ok {
		return false
	}

	if function.TemplateIndex != otherFunc.TemplateIndex {
		return false
	}

	if !slices.Equal(function.Captures, otherFunc.Captures) {
		return false
	}

	return true
}
