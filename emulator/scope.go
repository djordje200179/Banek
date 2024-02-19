package emulator

import (
	"banek/bytecode"
	"banek/runtime"
	"sync"
	"unsafe"
)

var scopeVarsPools = [...]sync.Pool{
	{New: func() any { return (*runtime.Obj)(nil) }},
	{New: func() any { return &(new([1]runtime.Obj)[0]) }},
	{New: func() any { return &(new([2]runtime.Obj)[0]) }},
	{New: func() any { return &(new([3]runtime.Obj)[0]) }},
}

func newScopeVars(size int) []runtime.Obj {
	if size >= len(scopeVarsPools) {
		return make([]runtime.Obj, size)
	}

	arr := scopeVarsPools[size].Get().(*runtime.Obj)
	slice := unsafe.Slice(arr, size)

	return slice
}

func freeScopeVars(slice []runtime.Obj) {
	if len(slice) >= len(scopeVarsPools) {
		return
	}

	clear(slice)
	scopeVarsPools[len(slice)].Put(unsafe.SliceData(slice))
}

type scope struct {
	pc   int
	vars []runtime.Obj

	function *bytecode.Func

	parent *scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}
