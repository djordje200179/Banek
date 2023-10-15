package expressions

import (
	"banek/ast"
	"strings"
)

type VariableAssignment struct {
	Variable Identifier
	Value    ast.Expression
}

func (assignment VariableAssignment) ExpressionNode() {}

func (assignment VariableAssignment) String() string {
	var sb strings.Builder

	sb.WriteString(assignment.Variable.String())
	sb.WriteString(" = ")
	sb.WriteString(assignment.Value.String())

	return sb.String()
}
