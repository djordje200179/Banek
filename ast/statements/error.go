package statements

type Error struct {
	Err error
}

func (statement Error) String() string {
	return statement.Err.Error()
}

func (statement Error) Error() string {
	return statement.Err.Error()
}

func (statement Error) HasSideEffects() bool {
	return true
}
