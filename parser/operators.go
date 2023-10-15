package parser

import "banek/tokens"

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

	tokens.LeftParenthesis: Call,

	tokens.Assign: Assignment,
}
