package bytecode

import (
	"banek/bytecode/instrs"
	"banek/runtime/builtins"
	"banek/runtime/ops"
	"fmt"
	"strconv"
	"strings"
)

type Code []byte

func (code Code) String() string {
	var sb strings.Builder

	for pc := 0; pc < len(code); {
		opcode, operands, width := instrs.ReadInstr(code[pc:])
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
			case instrs.OperandConst:
				sb.WriteByte('=')
				sb.WriteString(strconv.Itoa(operandValue))
			case instrs.OperandLiteral:
				sb.WriteString(strconv.Itoa(operandValue))
			case instrs.OperandFunc:
				sb.WriteByte('#')
				sb.WriteString(strconv.Itoa(operandValue))
			case instrs.OperandBinaryOp:
				sb.WriteString(ops.BinaryOperator(operandValue).String())
			case instrs.OperandUnaryOp:
				sb.WriteString(ops.UnaryOperator(operandValue).String())
			case instrs.OperandBuiltin:
				sb.WriteString(builtins.Funcs[operandValue].Name)
			}
		}

		sb.WriteByte('\n')

		pc += width
	}

	return sb.String()
}
