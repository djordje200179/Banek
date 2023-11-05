package instruction

import (
	"encoding/binary"
)

type OperandType byte

const (
	Constant OperandType = iota
	Literal
	Function

	InfixOperation
	PrefixOperation
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
		binary.PutVarint(code, int64(value))
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
		value, _ := binary.Varint(code)
		return int(value)
	}

	return 0
}
