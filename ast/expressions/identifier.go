package expressions

type Identifier string

func (identifier Identifier) ExpressionNode() {}
func (identifier Identifier) String() string  { return string(identifier) }
