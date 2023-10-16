package bytecode

import (
	"fmt"
	"strconv"
	"strings"
)

type Code []byte

func (code Code) String() string {
	var sb strings.Builder

	for pc := 0; pc < len(code); {
		operation, operands, width := ReadOperation(code[pc:])

		operandsStr := make([]string, len(operands))
		for i, operand := range operands {
			operandsStr[i] = "#" + strconv.Itoa(operand)
		}

		sb.WriteString(fmt.Sprintf("%04d %s %s\n", pc, operation, strings.Join(operandsStr, ",")))

		pc += width
	}

	return sb.String()
}
