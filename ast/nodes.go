package ast

import "fmt"

type Node interface {
	fmt.Stringer
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}
