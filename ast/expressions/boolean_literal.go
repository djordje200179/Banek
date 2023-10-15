package expressions

import (
	"strconv"
)

type BooleanLiteral bool

func (literal BooleanLiteral) ExpressionNode() {}

func (literal BooleanLiteral) String() string {
	return strconv.FormatBool(bool(literal))
}
