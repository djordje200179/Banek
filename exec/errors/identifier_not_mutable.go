package errors

type IdentifierNotMutableError struct {
	Identifier string
}

func (err IdentifierNotMutableError) Error() string {
	return "identifier not mutable: " + err.Identifier
}
