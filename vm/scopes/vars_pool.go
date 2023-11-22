package scopes

import (
	"banek/runtime/objs"
	"sync"
	"unsafe"
)

var scopeVarsPools = [...]sync.Pool{
	{New: func() any { return (*objs.Obj)(nil) }},
	{New: func() any { return &(new([1]objs.Obj)[0]) }},
	{New: func() any { return &(new([2]objs.Obj)[0]) }},
	{New: func() any { return &(new([3]objs.Obj)[0]) }},
}

func newScopeVars(size int) []objs.Obj {
	if size >= len(scopeVarsPools) {
		return make([]objs.Obj, size)
	}

	arr := scopeVarsPools[size].Get().(*objs.Obj)
	slice := unsafe.Slice(arr, size)

	return slice
}

func freeScopeVars(slice []objs.Obj) {
	if len(slice) >= len(scopeVarsPools) {
		return
	}

	clear(slice)
	scopeVarsPools[len(slice)].Put(unsafe.SliceData(slice))
}
