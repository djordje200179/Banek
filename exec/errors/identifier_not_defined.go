package errors

type ErrIdentifierNotDefined struct {
	Identifier string
}

func (err ErrIdentifierNotDefined) Error() string {
	return "identifier not defined: " + err.Identifier
}
