package objects

type Null struct{}

func (null Null) Type() ObjectType { return NullType }

func (null Null) Inspect() string { return "null" }
