package objects

type Return struct {
	Value Object
}

func (ret Return) Type() ObjectType { return ReturnType }

func (ret Return) String() string { return ret.Value.String() }
