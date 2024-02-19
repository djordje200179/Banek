package analyzer

import (
	"banek/ast/exprs"
	"banek/ast/stmts"
	"banek/tokens"
	"errors"
)

type UndefinedIdentError exprs.Ident

func (e UndefinedIdentError) Error() string { return "undefined identifier: " + e.String() }

type RedeclaredIdentError exprs.Ident

func (e RedeclaredIdentError) Error() string { return "redeclared identifier: " + e.String() }

var ErrReturnOutsideFunc = errors.New("return statement outside of function")

type InvalidAssignmentError stmts.Assignment

func (e InvalidAssignmentError) Error() string {
	return "invalid assignment: " + stmts.Assignment(e).String()
}

type InvalidOpError tokens.Type

func (e InvalidOpError) Error() string { return "invalid operator: " + tokens.Type(e).String() }
