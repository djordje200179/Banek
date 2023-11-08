package environments

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/objects"
	"strings"
)

type Func struct {
	Params []expressions.Identifier
	Body   ast.Statement

	Env Env
}

func (function *Func) Type() objects.Type    { return objects.TypeFunction }
func (function *Func) Clone() objects.Object { return function }

func (function *Func) String() string {
	var sb strings.Builder

	sb.WriteString("fn(")
	for i, param := range function.Params {
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
