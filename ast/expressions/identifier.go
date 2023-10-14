package expressions

type Identifier struct {
	Name string
}

func (identifier Identifier) ExpressionNode() {}

func (identifier Identifier) String() string {
	return identifier.Name
}
