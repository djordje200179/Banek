package bytecode

import (
	"banek/exec/objects"
)

type Executable struct {
	Code      Code
	Constants []objects.Object
}
