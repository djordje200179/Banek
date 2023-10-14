package expressions

import (
	"strconv"
)

type BooleanLiteral struct {
	Value bool
}

func (literal BooleanLiteral) ExpressionNode() {}

func (literal BooleanLiteral) String() string {
	return strconv.FormatBool(literal.Value)
}
