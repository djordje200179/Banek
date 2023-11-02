package ast

import "fmt"

type Statement interface {
	fmt.Stringer

	HasSideEffects() bool
}

type ErrUnknownStatement struct {
	Statement Statement
}

func (err ErrUnknownStatement) Error() string {
	return "unknown statement: " + err.Statement.String()
}
