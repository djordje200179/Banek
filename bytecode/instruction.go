package bytecode

import "encoding/binary"

func MakeInstruction(op Operation, operandValues ...int) Code {
	opInfo := op.Info()

	instruction := make(Code, opInfo.Size())
	instruction[0] = byte(op)

	offset := 1
	for i, operandValue := range operandValues {
		width := opInfo.Operands[i].Width

		copy(instruction[offset:], MakeOperand(operandValue, width))

		offset += width
	}
	return instruction
}

func ReadInstruction(instruction Code) (Operation, []int, int) {
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

func MakeOperand(value int, width int) Code {
	code := make(Code, width)
	switch width {
	case 2:
		binary.LittleEndian.PutUint16(code, uint16(value))
	default:
		return nil
	}

	return code
}

func ReadOperandValue(code Code, width int) int {
	switch width {
	case 2:
		return int(binary.LittleEndian.Uint16(code))
	}

	return 0
}
