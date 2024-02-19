package bytecode

import (
	"banek/bytecode/instrs"
	"strconv"
	"strings"
)

type Capture struct {
	Index int
	Level int
}

type FuncTemplate struct {
	Name      string
	NumParams int

	NumLocals int
	Code      instrs.Code

	IsCaptured bool

	Captures []Capture
}

func (t *FuncTemplate) IsClosure() bool { return len(t.Captures) > 0 }

func (t *FuncTemplate) String() string {
	var sb strings.Builder

	if t.Name != "" {
		sb.WriteString(t.Name)
	} else {
		sb.WriteString("<anonymous>")
	}

	sb.WriteByte('(')
	sb.WriteString(strconv.Itoa(t.NumParams))
	sb.WriteString(" params)\n")

	sb.WriteString(t.Code.String())

	return sb.String()
}
