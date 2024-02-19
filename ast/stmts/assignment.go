package stmts

import (
	"banek/ast"
	"banek/tokens"
)

type Assignment struct {
	Var, Value ast.Expr
}

func (stmt Assignment) String() string { return stmt.Var.String() + "=" + stmt.Value.String() }
func (stmt Assignment) StmtNode()      {}
func (stmt Assignment) DesStmtNode()   {}

type CompoundAssignment struct {
	Var, Value ast.Expr

	Operator tokens.Type
}

func (stmt CompoundAssignment) String() string {
	return stmt.Var.String() + stmt.Operator.String() + stmt.Value.String()
}

func (stmt CompoundAssignment) StmtNode()    {}
func (stmt CompoundAssignment) DesStmtNode() {}
