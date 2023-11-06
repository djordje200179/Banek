package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
	"strconv"
	"sync"
)

type Scope struct {
	code bytecode.Code
	pc   int

	vars []objects.Object

	function *bytecode.Func

	parent *Scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &Scope{}
	},
}

type ErrVarOutOfScope struct {
	Index int
}

func (err ErrVarOutOfScope) Error() string {
	return "variable out of Scope: " + strconv.Itoa(err.Index)
}

func (vm *vm) getGlobal(index int) (objects.Object, error) {
	if index >= len(vm.globalScope.vars) {
		return nil, ErrVarOutOfScope{index}
	}

	return vm.globalScope.vars[index], nil
}

func (vm *vm) setGlobal(index int, value objects.Object) error {
	if index >= len(vm.globalScope.vars) {
		return ErrVarOutOfScope{index}
	}

	vm.globalScope.vars[index] = value

	return nil
}

func (scope *Scope) getLocal(index int) (objects.Object, error) {
	if index >= len(scope.vars) {
		return nil, ErrVarOutOfScope{index}
	}

	return scope.vars[index], nil
}

func (scope *Scope) setLocal(index int, value objects.Object) error {
	if index >= len(scope.vars) {
		return ErrVarOutOfScope{index}
	}

	scope.vars[index] = value

	return nil
}

func (scope *Scope) getCaptured(index int) (objects.Object, error) {
	if index >= len(scope.function.Captures) {
		return nil, ErrVarOutOfScope{index}
	}

	return *scope.function.Captures[index], nil
}

func (scope *Scope) setCaptured(index int, value objects.Object) error {
	if index >= len(scope.function.Captures) {
		return ErrVarOutOfScope{index}
	}

	*scope.function.Captures[index] = value

	return nil
}

func (scope *Scope) hasCode() bool {
	return scope.pc < len(scope.code)
}

func (scope *Scope) readOpcode() instructions.Opcode {
	opcode := instructions.Opcode(scope.code[scope.pc])

	scope.pc++

	return opcode
}

func (scope *Scope) readOperand(width int) int {
	rawOperand := scope.code[scope.pc : scope.pc+width]
	operand := instructions.ReadOperandValue(rawOperand, width)

	scope.pc += width

	return operand
}

func (scope *Scope) movePC(offset int) {
	scope.pc += offset
}

func (vm *vm) pushScope(function *bytecode.Func, args []objects.Object) {
	funcTemplate := vm.program.FuncsPool[function.TemplateIndex]

	if len(args) > len(funcTemplate.Params) {
		args = args[:len(funcTemplate.Params)]
	}

	funcScope := scopePool.Get().(*Scope)

	funcScope.code = funcTemplate.Code
	funcScope.pc = 0
	funcScope.parent = vm.currScope
	funcScope.function = function

	if funcTemplate.NumLocals > len(args) {
		funcScope.vars = make([]objects.Object, funcTemplate.NumLocals)
		copy(funcScope.vars, args)
		for i := len(args); i < len(funcScope.vars); i++ {
			funcScope.vars[i] = objects.Undefined{}
		}
	} else {
		funcScope.vars = args
	}

	vm.currScope = funcScope
}

func (vm *vm) popScope() {
	removedScope := vm.currScope

	vm.currScope = vm.currScope.parent

	scopePool.Put(removedScope)
}
