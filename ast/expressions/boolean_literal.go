package expressions

import (
	"strconv"
)

type BooleanLiteral bool

func (literal BooleanLiteral) String() string {
	return strconv.FormatBool(bool(literal))
}

func (literal BooleanLiteral) IsConstant() bool { return true }
