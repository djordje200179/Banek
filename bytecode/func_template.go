package bytecode

import (
	"strings"
)

type Capture struct {
	Index int
	Level int
}

type FuncTemplate struct {
	Name string

	Code Code

	Params    []string
	NumLocals int

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
	sb.WriteString(strings.Join(template.Params, ", "))
	sb.WriteByte(')')
	sb.WriteByte('\n')

	sb.WriteString(template.Code.String())

	return sb.String()
}
