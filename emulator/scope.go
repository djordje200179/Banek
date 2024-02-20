package emulator

import (
	"banek/bytecode"
	"sync"
)

type scope struct {
	pc int
	bp int

	function *bytecode.Func

	parent *scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}
