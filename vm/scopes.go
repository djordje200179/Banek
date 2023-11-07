package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
	"strconv"
)

type scope struct {
	code bytecode.Code
	pc   int

	vars []objects.Object

	function *bytecode.Func

	parent *scope
}

type ErrVarOutOfScope struct {
	Index int
}

func (err ErrVarOutOfScope) Error() string {
	return "variable out of scope: " + strconv.Itoa(err.Index)
}

func (scope *scope) getLocal(index int) (objects.Object, error) {
	if index >= len(scope.vars) {
		return nil, ErrVarOutOfScope{index}
	}

	return scope.vars[index], nil
}

func (scope *scope) setLocal(index int, value objects.Object) error {
	if index >= len(scope.vars) {
		return ErrVarOutOfScope{index}
	}

	scope.vars[index] = value

	return nil
}

func (scope *scope) getCaptured(index int) (objects.Object, error) {
	if index >= len(scope.function.Captures) {
		return nil, ErrVarOutOfScope{index}
	}

	return *scope.function.Captures[index], nil
}

func (scope *scope) setCaptured(index int, value objects.Object) error {
	if index >= len(scope.function.Captures) {
		return ErrVarOutOfScope{index}
	}

	*scope.function.Captures[index] = value

	return nil
}

func (scope *scope) hasCode() bool {
	return scope.pc < len(scope.code)
}

func (scope *scope) readOpcode() instructions.Opcode {
	opcode := instructions.Opcode(scope.code[scope.pc])

	scope.pc++

	return opcode
}

func (scope *scope) readOperand(width int) int {
	rawOperand := scope.code[scope.pc : scope.pc+width]
	operand := instructions.ReadOperandValue(rawOperand, width)

	scope.pc += width

	return operand
}

func (scope *scope) movePC(offset int) {
	scope.pc += offset
}
