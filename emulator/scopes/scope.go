package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime"
	"sync"
)

type Scope struct {
	code instrs.Code

	pc   int
	vars []runtime.Obj

	function *bytecode.Func
	template *bytecode.FuncTemplate

	parent *Scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &Scope{}
	},
}
