package stmts

import (
	"banek/ast"
	"strings"
)

type For struct {
	Init ast.DesStmt
	Cond ast.Expr
	Post ast.DesStmt

	Body ast.Stmt
}

func (stmt For) String() string {
	var sb strings.Builder

	sb.WriteString("for")
	sb.WriteString(stmt.Init.String())
	sb.WriteString(";")
	sb.WriteString(stmt.Cond.String())
	sb.WriteString(";")
	sb.WriteString(stmt.Post.String())
	sb.WriteString(" do {\n")
	sb.WriteString(stmt.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

func (stmt For) StmtNode() {}
