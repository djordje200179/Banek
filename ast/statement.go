package ast

import "fmt"

type Statement interface {
	fmt.Stringer

	HasSideEffects() bool
}

type ErrUnknownStmt struct {
	Stmt Statement
}

func (err ErrUnknownStmt) Error() string {
	return "unknown statement: " + err.Stmt.String()
}
