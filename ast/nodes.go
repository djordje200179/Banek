package ast

import "fmt"

type Node interface {
	fmt.Stringer
}

type Statement interface {
	Node

	HasSideEffects() bool
}

type Expression interface {
	Node

	IsConstant() bool
}
