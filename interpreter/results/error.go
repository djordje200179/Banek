package results

type Error struct {
	Err error
}

func (err Error) String() string { return err.Err.Error() }
func (err Error) Error() string  { return err.Err.Error() }
