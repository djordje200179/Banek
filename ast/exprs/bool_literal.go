package exprs

import "strconv"

type BoolLiteral bool

func (b BoolLiteral) String() string { return strconv.FormatBool(bool(b)) }
func (b BoolLiteral) IsConst() bool  { return true }
