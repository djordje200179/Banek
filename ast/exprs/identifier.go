package exprs

type Identifier string

func (expr Identifier) String() string { return string(expr) }
func (expr Identifier) IsConst() bool  { return false }
