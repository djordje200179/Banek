package bytecode

import "banek/exec/objects"

type Function struct {
	TemplateIndex int

	Captures []*objects.Object
}

func (function Function) Type() string { return "function" }

func (function Function) String() string {
	return "<function>"
}
