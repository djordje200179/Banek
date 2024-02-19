package exprs

type UndefinedLiteral struct{}

func (_ UndefinedLiteral) String() string { return "undefined" }
func (_ UndefinedLiteral) IsConst() bool  { return true }
