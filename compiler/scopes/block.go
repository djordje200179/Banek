package scopes

import (
	"banek/bytecode/instruction"
	"fmt"
)

type Block struct {
	blocksCounter

	Index int

	Parent Scope
}

func (scope *Block) AddVar(name string, mutable bool) (int, error) {
	newName := fmt.Sprintf("%d#%s", scope.Index, name)

	return scope.Parent.AddVar(newName, mutable)
}

func (scope *Block) GetVar(name string) (Var, int) {
	newName := fmt.Sprintf("%d#%s", scope.Index, name)

	localVar, index := scope.Parent.GetVar(newName)
	if index != -1 {
		return localVar, index
	}

	return scope.Parent.GetVar(name)
}

func (scope *Block) EmitInstr(op instruction.Operation, operands ...int) {
	scope.Parent.EmitInstr(op, operands...)
}

func (scope *Block) PatchInstrOperand(addr int, operandIndex int, newValue int) {
	scope.Parent.PatchInstrOperand(addr, operandIndex, newValue)
}

func (scope *Block) CurrAddr() int {
	return scope.Parent.CurrAddr()
}

func (scope *Block) IsGlobal() bool {
	return scope.Parent.IsGlobal()
}
