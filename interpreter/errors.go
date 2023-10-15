package interpreter

import (
	"banek/ast"
	"banek/interpreter/objects"
	"banek/tokens"
	"fmt"
)

type UnknownStatementError struct {
	Statement ast.Statement
}

func (err UnknownStatementError) Error() string {
	return "unknown statement: " + err.Statement.String()
}

type IdentifierAlreadyDefinedError struct {
	Identifier string
}

func (err IdentifierAlreadyDefinedError) Error() string {
	return "identifier already defined: " + err.Identifier
}

type IdentifierNotDefinedError struct {
	Identifier string
}

func (err IdentifierNotDefinedError) Error() string {
	return "identifier not defined: " + err.Identifier
}

type IdentifierNotMutableError struct {
	Identifier string
}

func (err IdentifierNotMutableError) Error() string {
	return "identifier not mutable: " + err.Identifier
}

type UnknownOperatorError struct {
	Operator tokens.TokenType
}

func (err UnknownOperatorError) Error() string {
	return "unknown operator: " + err.Operator.String()
}

type UnknownExpressionError struct {
	Expression ast.Expression
}

func (err UnknownExpressionError) Error() string {
	return "unknown expression: " + err.Expression.String()
}

type InvalidOperandError struct {
	Operator string
	Operand  objects.Object
}

func (err InvalidOperandError) Error() string {
	return fmt.Sprintf("invalid operand: expected %s, got %s", err.Operator, err.Operand.Type())
}

type IndexOutOfBoundsError struct {
	Index int
	Size  int
}

func (err IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("index out of bounds: index %d, size %d", err.Index, err.Size)
}
