package envs

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/runtime/types"
	"slices"
	"strings"
)

type Func struct {
	Params []exprs.Identifier
	Body   ast.Stmt

	Env *Env
}

func (function *Func) Type() types.Type { return types.TypeFunc }
func (function *Func) Clone() types.Obj { return function }

func (function *Func) Equals(other types.Obj) bool {
	otherFunc, ok := other.(*Func)
	if !ok {
		return false
	}

	if !slices.Equal(function.Params, otherFunc.Params) {
		return false
	}

	if function.Body != otherFunc.Body {
		return false
	}

	if function.Env != otherFunc.Env {
		return false
	}

	return true
}

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
