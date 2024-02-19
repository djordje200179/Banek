package parser

import (
	"banek/ast"
	"banek/tokens"
	"fmt"
)

type UnexpectedTokenError struct {
	Expected, Got tokens.Type
}

func (err UnexpectedTokenError) Error() string {
	return fmt.Sprintf("expected %s, got %s instead", err.Expected.String(), err.Got.String())
}

type InvalidOperatorError tokens.Type

func (err InvalidOperatorError) Error() string {
	return "invalid operator: " + tokens.Type(err).String()
}

type InvalidExprStmtError struct {
	ast.Expr
}

func (err InvalidExprStmtError) Error() string {
	return "invalid expression statement: " + err.Expr.String()
}

type InvalidDesStmtError struct {
	ast.Stmt
}

func (err InvalidDesStmtError) Error() string {
	return "invalid designator statement: " + err.Stmt.String()
}
