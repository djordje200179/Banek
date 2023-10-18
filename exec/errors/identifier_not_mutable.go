package errors

type ErrIdentifierNotMutable struct {
	Identifier string
}

func (err ErrIdentifierNotMutable) Error() string {
	return "identifier not mutable: " + err.Identifier
}
