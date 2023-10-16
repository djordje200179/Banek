package instruction

func MakeInstruction(operation Operation, operandValues ...int) []byte {
	opInfo := operation.Info()

	instruction := make([]byte, opInfo.Size())
	instruction[0] = byte(operation)

	offset := 1
	for i, operandValue := range operandValues {
		width := opInfo.Operands[i].Width

		copy(instruction[offset:], MakeOperandValue(operandValue, width))

		offset += width
	}
	return instruction
}

func ReadInstruction(instruction []byte) (Operation, []int, int) {
	operation := Operation(instruction[0])
	opInfo := operation.Info()

	operandValues := make([]int, len(opInfo.Operands))

	offset := 1
	for i, operand := range opInfo.Operands {
		operandValues[i] = ReadOperandValue(instruction[offset:], operand.Width)
		offset += operand.Width
	}

	return operation, operandValues, offset
}
