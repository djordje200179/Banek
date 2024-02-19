package codegen

import "banek/bytecode/instrs"

type container struct {
	code instrs.Code

	level, index int
	vars         int

	previous *container
}

func (c *container) emitInstr(opcode instrs.Opcode, operands ...int) {
	instr := instrs.MakeInstr(opcode, operands...)

	newCode := make(instrs.Code, len(c.code)+len(instr))
	copy(newCode, c.code)
	copy(newCode[len(c.code):], instr)

	c.code = newCode
}

func (c *container) patchJumpOperand(addr int, operandIndex int) {
	op := instrs.Opcode(c.code[addr])
	opInfo := op.Info()

	instCode := c.code[addr : addr+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	offset := c.currAddr() - addr - opInfo.Size()
	copy(instCode[operandOffset:], instrs.MakeOperandValue(offset, operandWidth))
}

func (c *container) currAddr() int {
	return len(c.code)
}
