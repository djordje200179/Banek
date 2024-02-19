package parser

import (
	"banek/tokens"
)

type OperatorPrecedence byte

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

var infixOperatorPrecedences = map[tokens.Type]OperatorPrecedence{
	tokens.Equals:        Comparison,
	tokens.NotEquals:     Comparison,
	tokens.Less:          Comparison,
	tokens.Greater:       Comparison,
	tokens.LessEquals:    Comparison,
	tokens.GreaterEquals: Comparison,

	tokens.Plus:  Sum,
	tokens.Minus: Sum,

	tokens.Asterisk: Product,
	tokens.Slash:    Product,
	tokens.Percent:  Product,

	tokens.LParen:   Call,
	tokens.LBracket: Call,

	tokens.Assign:         Assignment,
	tokens.PlusAssign:     Assignment,
	tokens.MinusAssign:    Assignment,
	tokens.AsteriskAssign: Assignment,
	tokens.SlashAssign:    Assignment,
	tokens.PercentAssign:  Assignment,
	tokens.LArrow:         Assignment,
}
