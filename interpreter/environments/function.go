package environments

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/objects"
	"strings"
)

type Function struct {
	Parameters []expressions.Identifier
	Body       ast.Statement

	Env Environment
}

func (function *Function) Type() string          { return "function" }
func (function *Function) Clone() objects.Object { return function }

func (function *Function) String() string {
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
