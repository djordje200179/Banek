package vm

import (
	"banek/bytecode"
	"banek/runtime/objs"
	"sync"
	"unsafe"
)

type scope struct {
	code bytecode.Code

	pc   int
	vars []objs.Obj

	function *bytecode.Func

	parent *scope
}

type scopeStack struct {
	globalScope scope
	activeScope scope

	lastScope *scope
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

func (stack *scopeStack) getLocal(index int) objs.Obj {
	return stack.activeScope.vars[index]
}

func (stack *scopeStack) setLocal(index int, value objs.Obj) {
	stack.activeScope.vars[index] = value
}

func (stack *scopeStack) getCaptured(index int) objs.Obj {
	return *stack.activeScope.function.Captures[index]
}

func (stack *scopeStack) setCaptured(index int, value objs.Obj) {
	*stack.activeScope.function.Captures[index] = value
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

	return slice
}

func returnScopeVars(slice []objs.Obj) {
	if len(slice) >= len(scopeVarsPools) {
		return
	}

	for i := range slice {
		slice[i] = objs.Obj{}
	}
	scopeVarsPools[len(slice)].Put(unsafe.SliceData(slice))
}

func (stack *scopeStack) backupScope() {
	funcScope := scopePool.Get().(*scope)
	*funcScope = stack.activeScope
	stack.lastScope = funcScope
}

func (stack *scopeStack) restoreScope() {
	restoredScopeNode := stack.lastScope
	stack.lastScope = stack.lastScope.parent
	stack.activeScope = *restoredScopeNode

	restoredScopeNode.vars = nil
	restoredScopeNode.function = nil
	restoredScopeNode.parent = nil

	scopePool.Put(restoredScopeNode)
}
