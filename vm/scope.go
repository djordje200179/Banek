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
	index uint16
}

func (err ErrVariableOutOfScope) Error() string {
	return "variable out of scope: " + strconv.Itoa(int(err.index))
}

func (vm *vm) getGlobal(index uint16) (objects.Object, error) {
	globalScope := vm.scopeStack[0]

	if int(index) >= len(globalScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return globalScope.variables[index], nil
}

func (vm *vm) setGlobal(index uint16, value objects.Object) error {
	globalScope := vm.scopeStack[0]

	if int(index) >= len(globalScope.variables) {
		return ErrVariableOutOfScope{index}
	}

	globalScope.variables[index] = value

	return nil
}

func (vm *vm) getLocal(index uint16) (objects.Object, error) {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	if int(index) >= len(localScope.variables) {
		return nil, ErrVariableOutOfScope{index}
	}

	return localScope.variables[index], nil
}

func (vm *vm) setLocal(index uint16, value objects.Object) error {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	if int(index) >= len(localScope.variables) {
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

	return instruction.Operation(localScope.code[localScope.pc])
}

func (vm *vm) readCode() bytecode.Code {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	return localScope.code[localScope.pc+1:]
}

func (vm *vm) movePC(offset int) {
	localScope := vm.scopeStack[len(vm.scopeStack)-1]

	localScope.pc += offset
}
