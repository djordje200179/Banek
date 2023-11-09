package scopes

import "banek/bytecode/instrs"

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

	EmitInstr(opcode instrs.Opcode, operands ...int)
	PatchInstrOperand(addr int, operandIndex int, newValue int)
	CurrAddr() int

	IsGlobal() bool
	GetFunc() *Func

	MarkCaptured()

	NextBlockIndex() int
}
