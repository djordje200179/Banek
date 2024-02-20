package bytecode

import (
	"banek/bytecode/instrs"
	"strconv"
	"strings"
)

type Executable struct {
	Code instrs.Code

	StringPool []string
	FuncPool   []Func
}

func (e Executable) String() string {
	var sb strings.Builder

	replacePairs := make([]string, len(e.StringPool)*2)
	for i, str := range e.StringPool {
		replacePairs[i*2] = "$" + strconv.Itoa(i)
		replacePairs[i*2+1] = str
	}
	replacer := strings.NewReplacer(replacePairs...)

	sb.WriteString("Code:\n")
	sb.WriteString(e.Code.String())

	sb.WriteString("Functions:\n")
	for _, f := range e.FuncPool {
		sb.WriteString(replacer.Replace(f.String()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
