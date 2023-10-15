package objects

type None struct{}

func (none None) Type() ObjectType { return NoneType }

func (none None) String() string { return "--none--" }
