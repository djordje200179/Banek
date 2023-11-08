package vm

import (
	"banek/bytecode"
	"banek/exec/objects"
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

func (stack *scopeStack) getGlobal(index int) (objects.Object, error) {
	if index >= len(stack.globalScope.vars) {
		return nil, ErrVarOutOfScope{index}
	}

	return stack.globalScope.vars[index], nil
}

func (stack *scopeStack) setGlobal(index int, value objects.Object) error {
	if index >= len(stack.globalScope.vars) {
		return ErrVarOutOfScope{index}
	}

	stack.globalScope.vars[index] = value

	return nil
}

var scopeVarsPools = [...]sync.Pool{
	{New: func() any { return (*objects.Object)(nil) }},
	{New: func() any { return &(new([1]objects.Object)[0]) }},
	{New: func() any { return &(new([2]objects.Object)[0]) }},
	{New: func() any { return &(new([3]objects.Object)[0]) }},
	{New: func() any { return &(new([4]objects.Object)[0]) }},
}

func getScopeVars(size int) []objects.Object {
	arr := scopeVarsPools[size].Get().(*objects.Object)

	return unsafe.Slice(arr, size)
}

func returnScopeVars(arr []objects.Object) {
	scopeVarsPools[len(arr)].Put(unsafe.SliceData(arr))
}

func (stack *scopeStack) pushScope(code bytecode.Code, varsNum int, function *bytecode.Func, funcTemplate bytecode.FuncTemplate) *scope {
	funcScope := scopePool.Get().(*scope)

	funcScope.code = code
	funcScope.pc = 0
	funcScope.parent = stack.currScope
	funcScope.function = function
	funcScope.funcTemplate = funcTemplate

	if varsNum < len(scopeVarsPools) {
		funcScope.vars = getScopeVars(varsNum)
	} else {
		funcScope.vars = make([]objects.Object, varsNum)
		for i := range funcScope.vars {
			funcScope.vars[i] = objects.Undefined{}
		}
	}

	stack.currScope = funcScope

	return funcScope
}

func (stack *scopeStack) popScope() {
	removedScope := stack.currScope

	stack.currScope = stack.currScope.parent

	if removedScope.funcTemplate.IsCaptured {
		returnScopeVars(removedScope.vars)
	}

	removedScope.vars = nil
	scopePool.Put(removedScope)
}
