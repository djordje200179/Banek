package expressions

import "banek/exec/objects"

type ConstLiteral struct {
	Value objects.Object
}

func (expr ConstLiteral) String() string {
	return expr.Value.String()
}

func (expr ConstLiteral) IsConst() bool {
	return true
}
