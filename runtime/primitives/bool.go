package primitives

import (
	"banek/runtime"
	"strconv"
)

type Bool bool

func (b Bool) String() string                { return strconv.FormatBool(bool(b)) }
func (b Bool) Truthy() bool                  { return bool(b) }
func (b Bool) Clone() runtime.Obj            { return b }
func (b Bool) Equals(other runtime.Obj) bool { return b == other }

func (b Bool) Not() (runtime.Obj, bool) { return !b, true }
