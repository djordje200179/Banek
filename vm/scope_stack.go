package vm

import (
	"banek/bytecode"
	"banek/runtime/objs"
	"sync"
	"unsafe"
)

type scopeStack struct {
	globalScope scope
	*scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}

func (stack *scopeStack) getGlobal(index int) objs.Obj {
	return stack.globalScope.vars[index]
}

func (stack *scopeStack) setGlobal(index int, value objs.Obj) {
	stack.globalScope.vars[index] = value
}

var scopeVarsPools = [...]sync.Pool{
	{New: func() any { return (*objs.Obj)(nil) }},
	{New: func() any { return &(new([1]objs.Obj)[0]) }},
	{New: func() any { return &(new([2]objs.Obj)[0]) }},
	{New: func() any { return &(new([3]objs.Obj)[0]) }},
	{New: func() any { return &(new([4]objs.Obj)[0]) }},
}

func getScopeVars(size int) []objs.Obj {
	if size >= len(scopeVarsPools) {
		return make([]objs.Obj, size)
	}

	arr := scopeVarsPools[size].Get().(*objs.Obj)
	slice := unsafe.Slice(arr, size)

	for i := range slice {
		slice[i] = objs.Obj{}
	}

	return slice
}

func returnScopeVars(arr []objs.Obj) {
	if len(arr) >= len(scopeVarsPools) {
		return
	}

	scopeVarsPools[len(arr)].Put(unsafe.SliceData(arr))
}

func (stack *scopeStack) pushScope(code bytecode.Code, varsNum int, function *bytecode.Func) *scope {
	funcScope := scopePool.Get().(*scope)
	*funcScope = scope{
		code:     code,
		vars:     getScopeVars(varsNum),
		function: function,
		parent:   stack.scope,
	}
	stack.scope = funcScope

	return funcScope
}

func (stack *scopeStack) popScope(canFreeVars bool) {
	removedScope := stack.scope
	stack.scope = removedScope.parent

	if canFreeVars {
		returnScopeVars(removedScope.vars)
	}

	removedScope.vars = nil
	scopePool.Put(removedScope)
}
