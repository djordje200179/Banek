package expressions

import (
	"strconv"
)

type IntegerLiteral int

func (literal IntegerLiteral) String() string { return strconv.Itoa(int(literal)) }

func (literal IntegerLiteral) IsConstant() bool { return true }
