package stmts

import "banek/ast/exprs"

type FuncCall exprs.FuncCall

func (stmt FuncCall) String() string { return exprs.FuncCall(stmt).String() }
func (stmt FuncCall) StmtNode()      {}
