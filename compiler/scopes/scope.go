package scopes

import "banek/bytecode/instructions"

type Var struct {
	Name string

	Mutable bool
}

type blocksCounter int

func (counter *blocksCounter) NextBlockIndex() int {
	*counter++

	return int(*counter)
}

type Scope interface {
	AddVar(name string, mutable bool) (int, error)
	GetVar(name string) (Var, int)

	EmitInstr(opcode instructions.Opcode, operands ...int)
	PatchInstrOperand(addr int, operandIndex int, newValue int)
	CurrAddr() int

	IsGlobal() bool

	NextBlockIndex() int
}
