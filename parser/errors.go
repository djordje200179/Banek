package parser

import (
	"banek/ast"
	"banek/tokens"
	"strings"
)

type UnexpectedTokenError struct {
	Expected, Got tokens.Type
}

func (err UnexpectedTokenError) Error() string {
	var sb strings.Builder

	sb.WriteString("expected ")
	sb.WriteString(err.Expected.String())
	sb.WriteString(", got ")
	sb.WriteString(err.Got.String())
	sb.WriteString(" instead")

	return sb.String()
}

type InvalidTokenError struct {
	Got tokens.Type
}

func (err InvalidTokenError) Error() string {
	var sb strings.Builder

	sb.WriteString("invalid unary operator ")
	sb.WriteString(err.Got.String())

	return sb.String()
}

type InvalidExprStmtError struct {
	Expr ast.Expr
}

func (err InvalidExprStmtError) Error() string {
	return "invalid expression statement: " + err.Expr.String()
}
