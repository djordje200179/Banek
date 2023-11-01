package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
	"strconv"
	"sync"
)

type scope struct {
	code bytecode.Code
	pc   int

	variables []objects.Object

	parent *scope
}

var scopePool = sync.Pool{
	New: func() interface{} {
		return &scope{}
	},
}

type ErrVariableOutOfScope struct {
	Index int
}

func (err ErrVariableOutOfScope) Error() string {
	return "variable out of scope: " + strconv.Itoa(err.Index)
}

func (vm *vm) getGlobal(index int) (objects.Object, error) {
	if index >= len(vm.globalScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return vm.globalScope.variables[index], nil
}

func (vm *vm) setGlobal(index int, value objects.Object) error {
	if index >= len(vm.globalScope.variables) {
		return ErrVariableOutOfScope{index}
	}

	vm.globalScope.variables[index] = value

	return nil
}

func (vm *vm) getLocal(index int) (objects.Object, error) {
	if index >= len(vm.currentScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return vm.currentScope.variables[index], nil
}

func (vm *vm) setLocal(index int, value objects.Object) error {
	if index >= len(vm.currentScope.variables) {
		return ErrVariableOutOfScope{index}
	}

	vm.currentScope.variables[index] = value

	return nil
}

func (vm *vm) hasCode() bool {
	return vm.currentScope.pc < len(vm.currentScope.code)
}

func (vm *vm) readOperation() instruction.Operation {
	operation := instruction.Operation(vm.currentScope.code[vm.currentScope.pc])

	vm.currentScope.pc++

	return operation
}

func (vm *vm) readOperand(width int) int {
	rawOperand := vm.currentScope.code[vm.currentScope.pc : vm.currentScope.pc+width]
	operand := instruction.ReadOperandValue(rawOperand, width)

	vm.currentScope.pc += width

	return operand
}

func (vm *vm) movePC(offset int) {
	vm.currentScope.pc += offset
}
