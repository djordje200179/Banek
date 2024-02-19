package ast

import "fmt"

type Stmt interface {
	fmt.Stringer

	StmtNode()
}

type DesStmt interface {
	Stmt
	DesStmtNode()
}
