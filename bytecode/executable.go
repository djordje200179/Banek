package bytecode

import (
	"banek/exec/objects"
	"strconv"
	"strings"
)

type Executable struct {
	Code Code

	ConstantsPool []objects.Object
	FunctionsPool []FunctionTemplate

	NumGlobals int
}

func (executable Executable) String() string {
	var sb strings.Builder

	sb.WriteString("Code:\n")
	sb.WriteString(executable.Code.String())
	sb.WriteByte('\n')

	sb.WriteString("Constants:\n")
	for i, constant := range executable.ConstantsPool {
		sb.WriteString(strconv.Itoa(i) + ": " + constant.String())
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')

	sb.WriteString("Functions:\n")
	for i, function := range executable.FunctionsPool {
		sb.WriteString(strconv.Itoa(i))
		sb.WriteByte(':')
		sb.WriteString(function.String())
		sb.WriteByte('\n')
	}

	return sb.String()
}
