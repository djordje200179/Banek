package bytecode

import (
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
	Code      Code

	IsCaptured bool

	Captures []Capture
}

func (template FuncTemplate) IsClosure() bool {
	return len(template.Captures) > 0
}

func (template FuncTemplate) String() string {
	var sb strings.Builder

	if template.Name != "" {
		sb.WriteString(template.Name)
	} else {
		sb.WriteString("<anonymous>")
	}

	sb.WriteByte('(')
	sb.WriteString(strconv.Itoa(template.NumParams))
	sb.WriteString(" params)\n")

	sb.WriteString(template.Code.String())

	return sb.String()
}
