package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/errors"
	"slices"
)

type Global struct {
	vars []Var

	code bytecode.Code

	blocksCounter
}

func (scope *Global) AddVar(name string, mutable bool) (int, error) {
	if slices.ContainsFunc(scope.vars, func(v Var) bool {
		return v.Name == name
	}) {
		return 0, errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	scope.vars = append(scope.vars, Var{
		Name:    name,
		Mutable: mutable,
	})

	return len(scope.vars) - 1, nil
}

func (scope *Global) GetVar(name string) (Var, int) {
	index := slices.IndexFunc(scope.vars, func(v Var) bool {
		return v.Name == name
	})
	if index == -1 {
		return Var{}, -1
	}

	return scope.vars[index], index
}

func (scope *Global) EmitInstr(opcode instrs.Opcode, operands ...int) {
	instr := instrs.MakeInstr(opcode, operands...)

	newCode := make(bytecode.Code, len(scope.code)+len(instr))
	copy(newCode, scope.code)
	copy(newCode[len(scope.code):], instr)

	scope.code = newCode
}

func (scope *Global) PatchInstrOperand(addr int, operandIndex int, newValue int) {
	op := instrs.Opcode(scope.code[addr])
	opInfo := op.Info()

	instCode := scope.code[addr : addr+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	copy(instCode[operandOffset:], instrs.MakeOperandValue(newValue, operandWidth))

}

func (scope *Global) CurrAddr() int {
	return len(scope.code)
}

func (scope *Global) IsGlobal() bool {
	return true
}

func (scope *Global) GetFunc() *Func {
	panic("unreachable")
}

func (scope *Global) MarkCaptured() {}

func (scope *Global) MakeExecutable() bytecode.Executable {
	return bytecode.Executable{
		Code: scope.code,

		NumGlobals: len(scope.vars),
	}
}
