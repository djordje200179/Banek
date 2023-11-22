package scopes

import (
	"banek/bytecode"
	"banek/runtime/objs"
	"sync"
)

type Scope struct {
	code bytecode.Code

	pc   int
	vars []objs.Obj

	function *bytecode.Func
	template *bytecode.FuncTemplate

	parent *Scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &Scope{}
	},
}
