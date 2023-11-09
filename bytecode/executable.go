package bytecode

import (
	"banek/runtime/types"
	"strconv"
	"strings"
)

type Executable struct {
	Code Code

	ConstsPool []types.Obj
	FuncsPool  []FuncTemplate

	NumGlobals int
}

func (executable Executable) String() string {
	var sb strings.Builder

	replacePairs := make([]string, len(executable.ConstsPool)*2)
	for i, constant := range executable.ConstsPool {
		replacePairs[i*2] = "=" + strconv.Itoa(i)

		if functionObject, ok := constant.(*Func); ok {
			replacePairs[i*2+1] = executable.FuncsPool[functionObject.TemplateIndex].Name
		} else {
			replacePairs[i*2+1] = constant.String()
		}
	}
	replacer := strings.NewReplacer(replacePairs...)

	sb.WriteString("Code:\n")
	sb.WriteString(replacer.Replace(executable.Code.String()))
	sb.WriteByte('\n')

	sb.WriteString("Functions:\n")
	for _, function := range executable.FuncsPool {
		sb.WriteString(replacer.Replace(function.String()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
