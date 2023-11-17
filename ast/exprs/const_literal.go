package exprs

import (
	"banek/runtime/objs"
)

type ConstLiteral struct {
	Value objs.Obj
}

func (expr ConstLiteral) String() string {
	return expr.Value.String()
}

func (expr ConstLiteral) IsConst() bool {
	return true
}
