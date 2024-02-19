package bytecode

import (
	"strconv"
	"strings"
)

type Executable struct {
	StringPool []string
	FuncPool   []FuncTemplate
}

func (e *Executable) String() string {
	var sb strings.Builder

	replacePairs := make([]string, len(e.StringPool)*2)
	for i, str := range e.StringPool {
		replacePairs[i*2] = "$" + strconv.Itoa(i)
		replacePairs[i*2+1] = str
	}
	replacer := strings.NewReplacer(replacePairs...)

	for _, function := range e.FuncPool {
		sb.WriteString(replacer.Replace(function.String()))
		sb.WriteByte('\n')
	}

	return sb.String()
}
