package bytecode

import (
	"banek/exec/objects"
	"slices"
	"strconv"
)

type Func struct {
	TemplateIndex int

	Captures []*objects.Object
}

func (function *Func) Type() objects.Type    { return objects.TypeFunction }
func (function *Func) Clone() objects.Object { return function }
func (function *Func) String() string        { return "func#" + strconv.Itoa(function.TemplateIndex) }

func (function *Func) Equals(other objects.Object) bool {
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
