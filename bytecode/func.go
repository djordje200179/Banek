package bytecode

import (
	"banek/runtime"
	"strconv"
)

type Func struct {
	TemplateIndex int

	Captures []*runtime.Obj
}

func (f *Func) String() string                { return "func #" + strconv.Itoa(f.TemplateIndex) }
func (f *Func) Truthy() bool                  { return true }
func (f *Func) Clone() runtime.Obj            { return f }
func (f *Func) Equals(other runtime.Obj) bool { return f == other }
