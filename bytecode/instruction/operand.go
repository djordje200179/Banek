package instruction

import (
	"encoding/binary"
)

type OperandType byte

const (
	Constant OperandType = iota
	Literal
	Function
)

type OperandInfo struct {
	Width int
	Type  OperandType
}

func MakeOperandValue(value int, width int) []byte {
	code := make([]byte, width)
	switch width {
	case 1:
		code[0] = byte(value)
	case 2:
		binary.LittleEndian.PutUint16(code, uint16(value))
	default:
		return nil
	}

	return code
}

func ReadOperandValue(code []byte, width int) int {
	switch width {
	case 1:
		return int(code[0])
	case 2:
		return int(binary.LittleEndian.Uint16(code))
	}

	return 0
}
