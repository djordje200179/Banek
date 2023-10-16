package bytecode

import "encoding/binary"

func MakeInstruction(op Operation, operands ...int) Code {
	opInfo := op.Info()
	instructionLen := 1
	for _, w := range opInfo.OperandWidths {
		instructionLen += w
	}
	instruction := make(Code, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, operand := range operands {
		width := opInfo.OperandWidths[i]
		switch width {
		case 2:
			binary.LittleEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}
	return instruction
}

func ReadOperation(instruction Code) (Operation, []int, int) {
	operation := Operation(instruction[0])
	opInfo := operation.Info()

	operands := make([]int, len(opInfo.OperandWidths))

	offset := 1
	for i, width := range opInfo.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.LittleEndian.Uint16(instruction[offset:]))
		}
		offset += width
	}

	return operation, operands, offset
}
