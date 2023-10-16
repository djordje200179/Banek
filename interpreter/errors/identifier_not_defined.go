package errors

type IdentifierNotDefinedError struct {
	Identifier string
}

func (err IdentifierNotDefinedError) Error() string {
	return "identifier not defined: " + err.Identifier
}
