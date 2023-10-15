package expressions

import (
	"strconv"
)

type IntegerLiteral int

func (literal IntegerLiteral) ExpressionNode() {}
func (literal IntegerLiteral) String() string  { return strconv.Itoa(int(literal)) }
