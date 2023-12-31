package scopes

import (
	"banek/bytecode/instrs"
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

func (scope *Block) EmitInstr(opcode instrs.Opcode, operands ...int) {
	scope.Parent.EmitInstr(opcode, operands...)
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

func (scope *Block) GetFunc() *Func {
	block := scope.Parent
	for {
		switch scope := block.(type) {
		case *Func:
			return scope
		case *Block:
			block = scope.Parent
		default:
			panic("unreachable")
		}
	}
}

func (scope *Block) MarkCaptured() {
	scope.Parent.MarkCaptured()
}
