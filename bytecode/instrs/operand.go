package instrs

type OperandType byte

const (
	OperandLiteral OperandType = iota
	OperandString
	OperandFunc
	OperandBuiltin
	OperandOffset
)

type OperandInfo struct {
	Width int
	Type  OperandType
}

func MakeOperandValue(value int, width int) []byte {
	operand := make([]byte, width)

	for i := range width {
		operand[i] = byte(value >> uint(8*i))
	}

	return operand
}

func ReadOperandValue(operand []byte) int {
	value := 0
	for i := range operand {
		value |= int(operand[i]) << uint(8*i)
	}

	if operand[len(operand)-1]&0x80 != 0 {
		value |= ^0 << uint(8*len(operand))
	}

	return value
}
