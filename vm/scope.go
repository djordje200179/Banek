package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
	"strconv"
	"sync"
)

type scope struct {
	code bytecode.Code
	pc   int

	vars []objects.Object

	function *bytecode.Func

	parent *scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}

type ErrVarOutOfScope struct {
	Index int
}

func (err ErrVarOutOfScope) Error() string {
	return "variable out of scope: " + strconv.Itoa(err.Index)
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

func (vm *vm) getLocal(index int) (objects.Object, error) {
	if index >= len(vm.currScope.vars) {
		return nil, ErrVarOutOfScope{index}
	}

	return vm.currScope.vars[index], nil
}

func (vm *vm) setLocal(index int, value objects.Object) error {
	if index >= len(vm.currScope.vars) {
		return ErrVarOutOfScope{index}
	}

	vm.currScope.vars[index] = value

	return nil
}

func (vm *vm) getCaptured(index int) (objects.Object, error) {
	if index >= len(vm.currScope.function.Captures) {
		return nil, ErrVarOutOfScope{index}
	}

	return *vm.currScope.function.Captures[index], nil
}

func (vm *vm) setCaptured(index int, value objects.Object) error {
	if index >= len(vm.currScope.function.Captures) {
		return ErrVarOutOfScope{index}
	}

	*vm.currScope.function.Captures[index] = value

	return nil
}

func (vm *vm) hasCode() bool {
	return vm.currScope.pc < len(vm.currScope.code)
}

func (vm *vm) readOpcode() instructions.Opcode {
	opcode := instructions.Opcode(vm.currScope.code[vm.currScope.pc])

	vm.currScope.pc++

	return opcode
}

func (vm *vm) readOperand(width int) int {
	rawOperand := vm.currScope.code[vm.currScope.pc : vm.currScope.pc+width]
	operand := instructions.ReadOperandValue(rawOperand, width)

	vm.currScope.pc += width

	return operand
}

func (vm *vm) movePC(offset int) {
	vm.currScope.pc += offset
}
