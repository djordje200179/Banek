package exprs

import "strconv"

type IntLiteral int

func (i IntLiteral) String() string { return strconv.Itoa(int(i)) }
func (i IntLiteral) IsConst() bool  { return true }
