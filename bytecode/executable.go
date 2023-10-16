package bytecode

import (
	"banek/exec/objects"
)

type Executable struct {
	Code      Code
	Constants []objects.Object
}

func (executable *Executable) String() string {
	return executable.Code.String()
	// TODO: Replace constant indexes with constant values
}
