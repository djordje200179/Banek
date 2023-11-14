package vm

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/types"
)

type scope struct {
	code bytecode.Code
	pc   int

	vars []types.Obj

	function     *bytecode.Func
	funcTemplate bytecode.FuncTemplate

	parent *scope
}

func (scope *scope) getLocal(index int) types.Obj {
	return scope.vars[index]
}

func (scope *scope) setLocal(index int, value types.Obj) {
	scope.vars[index] = value
}

func (scope *scope) getCaptured(index int) types.Obj {
	return *scope.function.Captures[index]
}

func (scope *scope) setCaptured(index int, value types.Obj) {
	*scope.function.Captures[index] = value
}

func (scope *scope) readOpcode() instrs.Opcode {
	opcode := instrs.Opcode(scope.code[scope.pc])

	scope.pc++

	return opcode
}

func (scope *scope) readOperand(width int) int {
	rawOperand := scope.code[scope.pc : scope.pc+width]
	operand := instrs.ReadOperandValue(rawOperand, width)

	scope.pc += width

	return operand
}

func (scope *scope) movePC(offset int) {
	scope.pc += offset
}
