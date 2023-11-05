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

	replacePairs := make([]string, len(executable.ConstantsPool)*2)
	for i, constant := range executable.ConstantsPool {
		replacePairs[i*2] = "=" + strconv.Itoa(i)

		if functionObject, ok := constant.(*Function); ok {
			replacePairs[i*2+1] = executable.FunctionsPool[functionObject.TemplateIndex].Name
		} else {
			replacePairs[i*2+1] = constant.String()
		}
	}
	replacer := strings.NewReplacer(replacePairs...)

	sb.WriteString("Code:\n")
	sb.WriteString(replacer.Replace(executable.Code.String()))
	sb.WriteByte('\n')

	sb.WriteString("Functions:\n")
	for _, function := range executable.FunctionsPool {
		sb.WriteString(replacer.Replace(function.String()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
