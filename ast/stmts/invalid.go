package stmts

type Invalid struct {
	Err error
}

func (stmt Invalid) String() string {
	return stmt.Err.Error()
}

func (stmt Invalid) Error() string {
	return stmt.Err.Error()
}

func (stmt Invalid) HasSideEffects() bool {
	return true
}
