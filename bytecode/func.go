package bytecode

import (
	"cmp"
	"fmt"
)

type Capture struct {
	Index int
	Level int
}

type Func struct {
	Name      string
	NumParams int

	NumLocals int
	Addr      int

	IsCaptured bool

	Captures []Capture
}

func (f *Func) IsClosure() bool { return len(f.Captures) > 0 }

func (f *Func) String() string {
	return fmt.Sprintf("%s(%d params), starts at %04d", cmp.Or(f.Name, "<anonymous>"), f.NumParams, f.Addr)
}
