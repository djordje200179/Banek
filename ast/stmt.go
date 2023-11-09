package ast

import "fmt"

type Stmt interface {
	fmt.Stringer

	HasSideEffects() bool
}

type ErrUnknownStmt struct {
	Stmt Stmt
}

func (err ErrUnknownStmt) Error() string {
	return "unknown statement: " + err.Stmt.String()
}
