package vm

import (
	"banek/bytecode"
	"banek/exec/objects"
	"sync"
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

func (stack *scopeStack) pushScope(code bytecode.Code, vars []objects.Object, function *bytecode.Func) {
	funcScope := scopePool.Get().(*scope)

	funcScope.code = code
	funcScope.pc = 0
	funcScope.parent = stack.currScope
	funcScope.function = function
	funcScope.vars = vars

	stack.currScope = funcScope
}

func (stack *scopeStack) popScope() {
	removedScope := stack.currScope

	stack.currScope = stack.currScope.parent

	returnObjectArray(removedScope.vars)
	removedScope.vars = nil

	scopePool.Put(removedScope)
}
