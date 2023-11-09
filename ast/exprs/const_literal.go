package exprs

import "banek/runtime/types"

type ConstLiteral struct {
	Value types.Obj
}

func (expr ConstLiteral) String() string {
	return expr.Value.String()
}

func (expr ConstLiteral) IsConst() bool {
	return true
}
