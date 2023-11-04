package parser

import (
	"banek/exec/operations"
	"banek/tokens"
)

type OperatorPrecedence int

const (
	_ OperatorPrecedence = iota
	Lowest
	Assignment
	Comparison
	Sum
	Product
	Prefix
	Call
)

var infixOperatorPrecedences = map[tokens.TokenType]OperatorPrecedence{
	tokens.Equals:              Comparison,
	tokens.NotEquals:           Comparison,
	tokens.LessThan:            Comparison,
	tokens.GreaterThan:         Comparison,
	tokens.LessThanOrEquals:    Comparison,
	tokens.GreaterThanOrEquals: Comparison,

	tokens.Plus:  Sum,
	tokens.Minus: Sum,

	tokens.Asterisk: Product,
	tokens.Slash:    Product,
	tokens.Modulo:   Product,
	tokens.Caret:    Product,

	tokens.LeftParenthesis: Call,
	tokens.LeftBracket:     Call,

	tokens.Assign:         Assignment,
	tokens.PlusAssign:     Assignment,
	tokens.MinusAssign:    Assignment,
	tokens.AsteriskAssign: Assignment,
	tokens.SlashAssign:    Assignment,
}

var infixOperations = map[tokens.TokenType]operations.InfixOperationType{
	tokens.Plus:     operations.InfixPlusOperation,
	tokens.Minus:    operations.InfixMinusOperation,
	tokens.Asterisk: operations.InfixAsteriskOperation,
	tokens.Slash:    operations.InfixSlashOperation,
	tokens.Modulo:   operations.InfixModuloOperation,
	tokens.Caret:    operations.InfixCaretOperation,

	tokens.Equals:              operations.InfixEqualsOperation,
	tokens.NotEquals:           operations.InfixNotEqualsOperation,
	tokens.LessThan:            operations.InfixLessThanOperation,
	tokens.GreaterThan:         operations.InfixGreaterThanOperation,
	tokens.LessThanOrEquals:    operations.InfixLessThanOrEqualsOperation,
	tokens.GreaterThanOrEquals: operations.InfixGreaterThanOrEqualsOperation,
}

var prefixOperations = map[tokens.TokenType]operations.PrefixOperationType{
	tokens.Minus: operations.PrefixMinusOperation,
	tokens.Bang:  operations.PrefixBangOperation,
}
