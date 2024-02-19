package analyzer

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/tokens"
	"errors"
)

type UndefinedIdentError exprs.Ident

func (e UndefinedIdentError) Error() string { return "undefined identifier: " + e.String() }

type RedeclaredIdentError exprs.Ident

func (e RedeclaredIdentError) Error() string { return "redeclared identifier: " + e.String() }

var ErrReturnOutsideFunc = errors.New("return statement outside of function")

type InvalidAssignmentError struct {
	ast.Stmt
}

func (e InvalidAssignmentError) Error() string {
	return "invalid assignment: " + e.Stmt.String()
}

type InvalidOpError tokens.Type

func (e InvalidOpError) Error() string { return "invalid operator: " + tokens.Type(e).String() }

type UninitializedImmutableVarError exprs.Ident

func (e UninitializedImmutableVarError) Error() string {
	return "uninitialized immutable variable: " + e.String()
}

type ImmutableVarAssignmentError exprs.Ident

func (e ImmutableVarAssignmentError) Error() string {
	return "assignment to immutable variable: " + e.String()
}
