package bytecode

import (
	"banek/bytecode/instruction"
	"fmt"
	"strconv"
	"strings"
)

type Code []byte

func (code Code) String() string {
	var sb strings.Builder

	for pc := 0; pc < len(code); {
		operation, operandValues, width := instruction.ReadInstruction(code[pc:])
		opInfo := operation.Info()

		operandsStr := make([]string, len(operandValues))
		for i, operandValue := range operandValues {
			operand := opInfo.Operands[i]

			switch operand.Type {
			case instruction.Constant:
				operandsStr[i] = "=" + strconv.Itoa(operandValue)
			case instruction.Literal:
				operandsStr[i] = strconv.Itoa(operandValue)
			case instruction.Function:
				operandsStr[i] = "#" + strconv.Itoa(operandValue)
			}
		}

		sb.WriteString(fmt.Sprintf("%04d %s %s\n", pc, operation, strings.Join(operandsStr, ",")))

		pc += width
	}

	return sb.String()
}
