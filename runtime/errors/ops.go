package errors

type ErrUnknownOperator struct {
	Operator string
}

func (err ErrUnknownOperator) Error() string {
	return "unknown operator: " + err.Operator
}
