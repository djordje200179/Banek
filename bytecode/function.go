package bytecode

import (
	"banek/exec/objects"
	"fmt"
)

type Function struct {
	TemplateIndex int

	Captures []*objects.Object
}

func (function *Function) Type() string          { return "function" }
func (function *Function) Clone() objects.Object { return function }

func (function *Function) String() string {
	return fmt.Sprintf("func#%d", function.TemplateIndex)
}
