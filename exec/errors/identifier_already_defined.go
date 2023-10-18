package errors

type ErrIdentifierAlreadyDefined struct {
	Identifier string
}

func (err ErrIdentifierAlreadyDefined) Error() string {
	return "identifier already defined: " + err.Identifier
}
