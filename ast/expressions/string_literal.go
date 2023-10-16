package expressions

type StringLiteral string

func (literal StringLiteral) String() string { return string(literal) }

func (literal StringLiteral) IsConstant() bool { return true }
