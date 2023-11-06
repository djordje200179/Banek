package bytecode

import (
	"banek/exec/objects"
	"fmt"
)

type Func struct {
	TemplateIndex int

	Captures []*objects.Object
}

func (function *Func) Type() string          { return "function" }
func (function *Func) Clone() objects.Object { return function }

func (function *Func) String() string {
	return fmt.Sprintf("func#%d", function.TemplateIndex)
}
