package expressions

import (
	"banek/ast"
	"strings"
)

type CollectionAccess struct {
	Collection ast.Expression
	Key        ast.Expression
}

func (expression CollectionAccess) String() string {
	var sb strings.Builder

	sb.WriteString(expression.Collection.String())
	sb.WriteByte('[')
	sb.WriteString(expression.Key.String())
	sb.WriteByte(']')

	return sb.String()
}

func (expression CollectionAccess) IsConstant() bool {
	return expression.Collection.IsConstant() && expression.Key.IsConstant()
}
