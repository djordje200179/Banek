package expressions

import (
	"banek/ast"
	"strings"
)

type CollectionIndex struct {
	Collection ast.Expression
	Index      ast.Expression
}

func (index CollectionIndex) ExpressionNode() {}
func (index CollectionIndex) String() string {
	var sb strings.Builder

	sb.WriteString(index.Collection.String())
	sb.WriteByte('[')
	sb.WriteString(index.Index.String())
	sb.WriteByte(']')

	return sb.String()
}
