package expressions

type StringLiteral string

func (literal StringLiteral) ExpressionNode() {}
func (literal StringLiteral) String() string  { return string(literal) }
