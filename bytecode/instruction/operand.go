package instruction

import (
	"encoding/binary"
)

type OperandType int

const (
	Constant OperandType = iota
	Literal
)

type OperandInfo struct {
	Width int
	Type  OperandType
}

var constantPoolOperand = OperandInfo{2, Constant}

func MakeOperandValue(value int, width int) []byte {
	code := make([]byte, width)
	switch width {
	case 2:
		binary.LittleEndian.PutUint16(code, uint16(value))
	default:
		return nil
	}

	return code
}

func ReadOperandValue(code []byte, width int) int {
	switch width {
	case 2:
		return int(binary.LittleEndian.Uint16(code))
	}

	return 0
}
