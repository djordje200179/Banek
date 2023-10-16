package expressions

type Identifier string

func (identifier Identifier) String() string { return string(identifier) }

func (identifier Identifier) IsConstant() bool { return false }
