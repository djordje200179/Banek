package ast

import "strings"

type Program struct {
	Statements []Statement
}

func (program Program) String() string {
	var sb strings.Builder

	for _, statement := range program.Statements {
		sb.WriteString(statement.String())
		sb.WriteByte('\n')
	}

	return sb.String()
}
