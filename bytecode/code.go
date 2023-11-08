package bytecode

import (
	"banek/bytecode/instructions"
	"banek/exec/operations"
	"fmt"
	"strconv"
	"strings"
)

type Code []byte

func (code Code) String() string {
	var sb strings.Builder

	for pc := 0; pc < len(code); {
		opcode, operands, width := instructions.ReadInstr(code[pc:])
		instrInfo := opcode.Info()

		sb.WriteString(fmt.Sprintf("%04d", pc))
		sb.WriteString(": ")
		sb.WriteString(opcode.String())

		for i, operandValue := range operands {
			if i > 0 {
				sb.WriteByte(',')
			} else {
				sb.WriteByte(' ')
			}

			operand := instrInfo.Operands[i]

			switch operand.Type {
			case instructions.OperandConst:
				sb.WriteByte('=')
				sb.WriteString(strconv.Itoa(operandValue))
			case instructions.OperandLiteral:
				sb.WriteString(strconv.Itoa(operandValue))
			case instructions.OperandFunc:
				sb.WriteByte('#')
				sb.WriteString(strconv.Itoa(operandValue))
			case instructions.OperandBinaryOp:
				sb.WriteString(operations.BinaryOperator(operandValue).String())
			case instructions.OperandUnaryOp:
				sb.WriteString(operations.UnaryOperator(operandValue).String())
			}
		}

		sb.WriteByte('\n')

		pc += width
	}

	return sb.String()
}
