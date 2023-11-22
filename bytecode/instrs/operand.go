package instrs

import (
	"unsafe"
)

type OperandType byte

const (
	OperandConst OperandType = iota
	OperandLiteral
	OperandFunc
	OperandBuiltin

	OperandBinaryOp
	OperandUnaryOp
)

type OperandInfo struct {
	Width int
	Type  OperandType
}

func MakeOperandValue(value int, width int) []byte {
	operand := make([]byte, width)
	ptr := unsafe.Pointer(unsafe.SliceData(operand))

	switch width {
	case 1:
		*(*int8)(ptr) = int8(value)
	case 2:
		*(*int16)(ptr) = int16(value)
	case 4:
		*(*int32)(ptr) = int32(value)
	case 8:
		*(*int64)(ptr) = int64(value)
	default:
		panic("invalid operand width")
	}

	return operand
}

func ReadOperandValue(operand []byte) int {
	ptr := unsafe.Pointer(unsafe.SliceData(operand))
	data := *(*int)(ptr)

	width := len(operand)
	width *= 8

	mask := ^uint(0) >> uint(64-width)
	data &= int(mask)

	return data
}
