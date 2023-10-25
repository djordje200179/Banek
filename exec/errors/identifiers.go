package errors

type ErrIdentifierNotDefined struct {
	Identifier string
}

func (err ErrIdentifierNotDefined) Error() string {
	return "identifier not defined: " + err.Identifier
}

type ErrIdentifierAlreadyDefined struct {
	Identifier string
}

func (err ErrIdentifierAlreadyDefined) Error() string {
	return "identifier already defined: " + err.Identifier
}

type ErrIdentifierNotMutable struct {
	Identifier string
}

func (err ErrIdentifierNotMutable) Error() string {
	return "identifier not mutable: " + err.Identifier
}
