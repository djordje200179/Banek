package parser

import (
	"banek/runtime/ops"
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

var binaryOps = map[tokens.Type]ops.BinaryOperator{
	tokens.Plus:     ops.BinaryPlus,
	tokens.Minus:    ops.BinaryMinus,
	tokens.Asterisk: ops.BinaryAsterisk,
	tokens.Slash:    ops.BinarySlash,
	tokens.Modulo:   ops.BinaryModulo,
	tokens.Caret:    ops.BinaryCaret,

	tokens.Equals:        ops.BinaryEquals,
	tokens.NotEquals:     ops.BinaryNotEquals,
	tokens.Less:          ops.BinaryLess,
	tokens.Greater:       ops.BinaryGreater,
	tokens.LessEquals:    ops.BinaryLessEquals,
	tokens.GreaterEquals: ops.BinaryGreaterEquals,
}

var unaryOps = map[tokens.Type]ops.UnaryOperator{
	tokens.Minus: ops.UnaryMinus,
	tokens.Bang:  ops.UnaryBang,
}
