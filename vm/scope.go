package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
	"strconv"
)

type scope struct {
	code bytecode.Code
	pc   int

	variables []objects.Object

	parent *scope
}

type ErrVariableOutOfScope struct {
	Index int
}

func (err ErrVariableOutOfScope) Error() string {
	return "variable out of scope: " + strconv.Itoa(err.Index)
}

func (vm *vm) getGlobal(index int) (objects.Object, error) {
	globalScope := vm.scopeStack[0]

	if index >= len(globalScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return globalScope.variables[index], nil
}

func (vm *vm) setGlobal(index int, value objects.Object) error {
	globalScope := vm.scopeStack[0]

	if index >= len(globalScope.variables) {
		return ErrVariableOutOfScope{index}
	}

	globalScope.variables[index] = value

	return nil
}

func (vm *vm) getLocal(index int) (objects.Object, error) {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	if index >= len(localScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return localScope.variables[index], nil
}

func (vm *vm) setLocal(index int, value objects.Object) error {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	if index >= len(localScope.variables) {
		return ErrVariableOutOfScope{index}
	}

	localScope.variables[index] = value

	return nil
}

func (vm *vm) hasCode() bool {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	return localScope.pc < len(localScope.code)
}

func (vm *vm) readOperation() instruction.Operation {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	code := localScope.code[localScope.pc]

	localScope.pc++

	return instruction.Operation(code)
}

func (vm *vm) readOperand(width int) int {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	rawOperand := localScope.code[localScope.pc : localScope.pc+width]
	operand := instruction.ReadOperandValue(rawOperand, width)

	localScope.pc += width

	return operand
}

func (vm *vm) movePC(offset int) {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	localScope.pc += offset
}
