package primitives

import "banek/runtime"

type Undefined struct{}

func (u Undefined) String() string                { return "undefined" }
func (u Undefined) Truthy() bool                  { return false }
func (u Undefined) Clone() runtime.Obj            { return u }
func (u Undefined) Equals(other runtime.Obj) bool { return u == other }
