package expressions

import "banek/exec/objects"

type ConstantLiteral struct {
	Value objects.Object
}

func (literal ConstantLiteral) String() string {
	return literal.Value.String()
}

func (literal ConstantLiteral) IsConstant() bool {
	return true
}
