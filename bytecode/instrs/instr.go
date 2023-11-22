package instrs

func MakeInstr(opcode Opcode, operands ...int) []byte {
	instrInfo := opcode.Info()

	instr := make([]byte, instrInfo.Size())
	instr[0] = byte(opcode)

	offset := 1
	for i, operand := range operands {
		width := instrInfo.Operands[i].Width

		copy(instr[offset:], MakeOperandValue(operand, width))

		offset += width
	}

	return instr
}

func ReadInstr(instr []byte) (Opcode, []int, int) {
	opcode := Opcode(instr[0])
	instrInfo := opcode.Info()

	operands := make([]int, len(instrInfo.Operands))

	offset := 1
	for i, operand := range instrInfo.Operands {
		operands[i] = ReadOperandValue(instr[offset : offset+operand.Width])
		offset += operand.Width
	}

	return opcode, operands, offset
}
