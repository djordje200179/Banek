package expressions

import (
	"banek/ast"
	"strings"
)

type ArrayLiteral []ast.Expression

func (literal ArrayLiteral) String() string {
	var sb strings.Builder

	elementsRepresentation := make([]string, len(literal))
	for i, element := range literal {
		elementsRepresentation[i] = element.String()
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elementsRepresentation, ", "))
	sb.WriteByte(']')

	return sb.String()
}

func (literal ArrayLiteral) IsConstant() bool {
	for _, element := range literal {
		if !element.IsConstant() {
			return false
		}
	}

	return true
}
