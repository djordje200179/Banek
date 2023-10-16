package bytecode

import (
	"banek/exec/objects"
	"strconv"
	"strings"
)

type Executable struct {
	Code Code

	ConstantsPool []objects.Object
}

func (executable *Executable) String() string {
	code := executable.Code.String()

	replacePairs := make([]string, len(executable.ConstantsPool)*2)
	for i, constant := range executable.ConstantsPool {
		replacePairs[i*2] = "#" + strconv.Itoa(i)
		replacePairs[i*2+1] = constant.String()
	}

	code = strings.NewReplacer(replacePairs...).Replace(code)

	return code
}
