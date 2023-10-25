package expressions

import (
	"banek/ast"
)

type UndefinedLiteral struct{}

func (literal UndefinedLiteral) String() string {
	return "undefined"
}

func (literal UndefinedLiteral) IsConstant() bool {
	return true
}

var UndefinedLiteralConst ast.Expression = UndefinedLiteral{}
