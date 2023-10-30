package bytecode

import (
	"strings"
)

type CaptureInfo struct {
	Index int
	Level int
}

type FunctionTemplate struct {
	Name string

	Code Code

	Parameters   []string
	NumLocals    int
	CapturesInfo []CaptureInfo
}

func (template FunctionTemplate) IsClosure() bool {
	return len(template.CapturesInfo) > 0
}

func (template FunctionTemplate) String() string {
	var sb strings.Builder

	if template.Name != "" {
		sb.WriteString(template.Name)
	} else {
		sb.WriteString("<anonymous>")
	}

	sb.WriteByte('(')
	sb.WriteString(strings.Join(template.Parameters, ", "))
	sb.WriteByte(')')
	sb.WriteByte('\n')

	sb.WriteString(template.Code.String())

	return sb.String()
}
