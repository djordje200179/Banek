package errors

type IdentifierAlreadyDefinedError struct {
	Identifier string
}

func (err IdentifierAlreadyDefinedError) Error() string {
	return "identifier already defined: " + err.Identifier
}
