package evaluator

type Error struct {
	Err error
}

func (err Error) String() string { return err.Err.Error() }
func (err Error) Error() string  { return err.Err.Error() }

type None struct{}

func (none None) String() string { return "--none--" }

type Return struct {
	Value Object
}

func (ret Return) String() string { return ret.Value.String() }
