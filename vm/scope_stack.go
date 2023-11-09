package vm

import (
	"banek/bytecode"
	"banek/runtime/objs"
	"banek/runtime/types"
	"sync"
	"unsafe"
)

type scopeStack struct {
	globalScope scope
	currScope   *scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}

func (stack *scopeStack) getGlobal(index int) types.Obj {
	return stack.globalScope.vars[index]
}

func (stack *scopeStack) setGlobal(index int, value types.Obj) {
	stack.globalScope.vars[index] = value
}

var scopeVarsPools = [...]sync.Pool{
	{New: func() any { return (*types.Obj)(nil) }},
	{New: func() any { return &(new([1]types.Obj)[0]) }},
	{New: func() any { return &(new([2]types.Obj)[0]) }},
	{New: func() any { return &(new([3]types.Obj)[0]) }},
	{New: func() any { return &(new([4]types.Obj)[0]) }},
}

func getScopeVars(size int) []types.Obj {
	if size >= len(scopeVarsPools) {
		return make([]types.Obj, size)
	}

	arr := scopeVarsPools[size].Get().(*types.Obj)

	return unsafe.Slice(arr, size)
}

func returnScopeVars(arr []types.Obj) {
	if len(arr) >= len(scopeVarsPools) {
		return
	}

	scopeVarsPools[len(arr)].Put(unsafe.SliceData(arr))
}

func (stack *scopeStack) pushScope(code bytecode.Code, varsNum int, function *bytecode.Func, funcTemplate bytecode.FuncTemplate) *scope {
	funcScope := scopePool.Get().(*scope)

	funcScope.code = code
	funcScope.pc = 0
	funcScope.parent = stack.currScope
	funcScope.function = function
	funcScope.funcTemplate = funcTemplate

	funcScope.vars = getScopeVars(varsNum)
	for i := range funcScope.vars {
		funcScope.vars[i] = objs.Undefined{}
	}

	stack.currScope = funcScope

	return funcScope
}

func (stack *scopeStack) popScope() {
	removedScope := stack.currScope

	stack.currScope = stack.currScope.parent

	if !removedScope.funcTemplate.IsCaptured {
		returnScopeVars(removedScope.vars)
	}

	removedScope.vars = nil
	scopePool.Put(removedScope)
}
