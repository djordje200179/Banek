package objects

import (
	"banek/ast"
	"banek/ast/expressions"
	"strings"
)

type Function struct {
	Parameters []expressions.Identifier
	Body       ast.Statement

	Env interface {
		Get(name string) (Object, error)
		Set(name string, value Object) error
	}
}

func (function Function) Type() string { return "function" }

func (function Function) String() string {
	var sb strings.Builder

	sb.WriteString("fn(")
	for i, param := range function.Parameters {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(function.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}
