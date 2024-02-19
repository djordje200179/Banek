package exprs

type StringLiteral string

func (s StringLiteral) String() string { return string(s) }
func (s StringLiteral) IsConst() bool  { return true }
