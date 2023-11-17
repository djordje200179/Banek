package envs

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/runtime/objs"
	"slices"
	"strings"
	"unsafe"
)

type Func struct {
	Params []exprs.Identifier
	Body   ast.Stmt

	Env *Env
}

func GetFunc(obj objs.Obj) *Func {
	return (*Func)(obj.PtrData)
}

func (function *Func) MakeObj() objs.Obj {
	ptrData := unsafe.Pointer(function)

	return objs.Obj{Tag: objs.TypeFunc, PtrData: ptrData}
}

func FuncEquals(first, second objs.Obj) bool {
	firstFunc := GetFunc(first)
	secondFunc := GetFunc(second)

	if !slices.Equal(firstFunc.Params, secondFunc.Params) {
		return false
	}

	if firstFunc.Body != secondFunc.Body {
		return false
	}

	if firstFunc.Env != secondFunc.Env {
		return false
	}

	return true
}

func FuncString(obj objs.Obj) string {
	function := GetFunc(obj)

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

func init() {
	objs.Config[objs.TypeFunc] = objs.TypeConfig{
		Stringer: FuncString,
		Equaler:  FuncEquals,
	}
}
