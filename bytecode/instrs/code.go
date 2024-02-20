package instrs

import (
	"fmt"
	"strconv"
	"strings"
)

type Code []byte

func (code Code) String() string {
	var sb strings.Builder

	for pc := 0; pc < len(code); {
		opcode, operands, width := ReadInstr(code[pc:])
		instrInfo := opcode.Info()

		sb.WriteString(fmt.Sprintf("%04d", pc))
		sb.WriteString(": ")
		sb.WriteString(opcode.String())

		for i, operandValue := range operands {
			if i > 0 {
				sb.WriteString(", ")
			} else {
				sb.WriteByte(' ')
			}

			operand := instrInfo.Operands[i]

			switch operand.Type {
			case OperandString:
				sb.WriteByte('$')
				sb.WriteString(strconv.Itoa(operandValue))
			case OperandLiteral:
				sb.WriteString(strconv.Itoa(operandValue))
			case OperandFunc:
				sb.WriteByte('#')
				sb.WriteString(strconv.Itoa(operandValue))
			case OperandOffset:
				sb.WriteString(strconv.Itoa(operandValue))
				sb.WriteString(" (=")
				sb.WriteString(strconv.Itoa(pc + width + operandValue))
				sb.WriteByte(')')
			default:
				panic("unreachable")
			}
		}

		sb.WriteByte('\n')

		pc += width
	}

	return sb.String()
}
