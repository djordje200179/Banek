package bytecode

import (
	"cmp"
	"fmt"
)

type Capture struct {
	Index int
	Level int
}

type FuncTemplate struct {
	Name      string
	NumParams int

	NumLocals int
	StartPC   int

	IsCaptured bool

	Captures []Capture
}

func (t *FuncTemplate) IsClosure() bool { return len(t.Captures) > 0 }

func (t *FuncTemplate) String() string {
	return fmt.Sprintf("%s(%d params), starts at %04d", cmp.Or(t.Name, "<anonymous>"), t.NumParams, t.StartPC)
}
