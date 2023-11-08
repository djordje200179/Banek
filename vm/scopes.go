package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
)

type scope struct {
	code bytecode.Code
	pc   int

	vars []objects.Object

	function     *bytecode.Func
	funcTemplate bytecode.FuncTemplate

	parent *scope
}

func (scope *scope) getLocal(index int) objects.Object {
	return scope.vars[index]
}

func (scope *scope) setLocal(index int, value objects.Object) {
	scope.vars[index] = value
}

func (scope *scope) getCaptured(index int) objects.Object {
	return *scope.function.Captures[index]
}

func (scope *scope) setCaptured(index int, value objects.Object) {
	*scope.function.Captures[index] = value
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
