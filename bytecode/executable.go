package bytecode

import (
	"banek/exec/objects"
)

type Executable struct {
	Code Code

	ConstantsPool []objects.Object
}

func (executable *Executable) String() string {
	return executable.Code.String()
	// TODO: Replace constant indexes with constant values
}
