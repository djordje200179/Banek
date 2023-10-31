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

	// TODO: Replace generic function object names with real names
	replacePairs := make([]string, (len(executable.ConstantsPool)+len(executable.FunctionsPool))*2)
	i := 0
	for _, constant := range executable.ConstantsPool {
		replacePairs[i*2] = "=" + strconv.Itoa(i)
		replacePairs[i*2+1] = constant.String()
		i++
	}
	for _, function := range executable.FunctionsPool {
		replacePairs[i*2] = "#" + strconv.Itoa(i)
		replacePairs[i*2+1] = function.String()
		i++
	}

	replacer := strings.NewReplacer(replacePairs...)

	sb.WriteString("Code:\n")
	sb.WriteString(replacer.Replace(executable.Code.String()))
	sb.WriteByte('\n')

	sb.WriteString("Functions:\n")
	for i, function := range executable.FunctionsPool {
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(": ")
		sb.WriteString(replacer.Replace(function.String()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
