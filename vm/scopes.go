package vm

import (
	"banek/bytecode"
	"banek/runtime/objs"
	"sync"
	"unsafe"
)

type scope struct {
	code bytecode.Code

	savedPC  int
	vars     []objs.Obj
	captures []*objs.Obj

	canFreeVars bool

	parent *scope
}

type scopeStack struct {
	globalScope scope
	activeScope *scope
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

func (stack *scopeStack) pushScope(funcTemplate *bytecode.FuncTemplate, function *bytecode.Func) *scope {
	funcScope := scopePool.Get().(*scope)
	*funcScope = scope{
		code:        funcTemplate.Code,
		vars:        getScopeVars(funcTemplate.NumLocals),
		captures:    function.Captures,
		canFreeVars: !funcTemplate.IsCaptured,
		parent:      stack.activeScope,
	}
	stack.activeScope = funcScope

	return funcScope
}

func (stack *scopeStack) popScope() {
	removedScope := stack.activeScope
	stack.activeScope = removedScope.parent

	if removedScope.canFreeVars {
		returnScopeVars(removedScope.vars)
	}

	removedScope.vars = nil
	scopePool.Put(removedScope)
}
