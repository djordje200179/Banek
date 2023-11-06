package parser

import (
	"banek/exec/operations"
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
	tokens.Modulo:   Product,
	tokens.Caret:    Product,

	tokens.LeftParen:   Call,
	tokens.LeftBracket: Call,

	tokens.Assign:         Assignment,
	tokens.PlusAssign:     Assignment,
	tokens.MinusAssign:    Assignment,
	tokens.AsteriskAssign: Assignment,
	tokens.SlashAssign:    Assignment,
}

var binaryOps = map[tokens.Type]operations.BinaryOperator{
	tokens.Plus:     operations.BinaryPlus,
	tokens.Minus:    operations.BinaryMinus,
	tokens.Asterisk: operations.BinaryAsterisk,
	tokens.Slash:    operations.BinarySlash,
	tokens.Modulo:   operations.BinaryModulo,
	tokens.Caret:    operations.BinaryCaret,

	tokens.Equals:        operations.BinaryEquals,
	tokens.NotEquals:     operations.BinaryNotEquals,
	tokens.Less:          operations.BinaryLess,
	tokens.Greater:       operations.BinaryGreater,
	tokens.LessEquals:    operations.BinaryLessEquals,
	tokens.GreaterEquals: operations.BinaryGreaterEquals,
}

var unaryOps = map[tokens.Type]operations.UnaryOperator{
	tokens.Minus: operations.UnaryMinus,
	tokens.Bang:  operations.UnaryBang,
}
