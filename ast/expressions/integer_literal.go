package expressions

import (
	"strconv"
)

type IntegerLiteral struct {
	Value int64
}

func (literal IntegerLiteral) ExpressionNode() {}

func (literal IntegerLiteral) String() string {
	return strconv.Itoa(int(literal.Value))
}
