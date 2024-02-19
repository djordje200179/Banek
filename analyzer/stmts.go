package analyzer

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/ast/stmts"
)

func (a *analyzer) analyzeStmt(stmt ast.Stmt) (ast.Stmt, error) {
	switch stmt := stmt.(type) {
	case stmts.Assignment:
		return a.analyzeAssignment(stmt)
	case stmts.CompoundAssignment:
		return a.analyzeCompoundAssignment(stmt)
	case stmts.Block:
		return a.analyzeStmtBlock(stmt)
	case stmts.FuncCall:
		expr, err := a.analyzeFuncCall(exprs.FuncCall(stmt))
		return stmts.FuncCall(expr), err
	case stmts.FuncDecl:
		return a.analyzeFuncDecl(stmt)
	case stmts.If:
		return a.analyzeIfStmt(stmt)
	case stmts.Return:
		return a.analyzeReturn(stmt)
	case stmts.VarDecl:
		return a.analyzeVarDecl(stmt)
	case stmts.While:
		return a.analyzeWhile(stmt)
	default:
		panic("unreachable")
	}
}

func (a *analyzer) analyzeStmtBlock(block stmts.Block) (stmts.Block, error) {
	a.symTable.OpenScope()
	defer a.symTable.CloseScope()
	for i, line := range block {
		line, err := a.analyzeStmt(line)
		if err != nil {
			return nil, err
		}
		block[i] = line
	}

	return block, nil
}

func (a *analyzer) analyzeWhile(stmt stmts.While) (ast.Stmt, error) {
	a.loopCnt++
	defer func() { a.loopCnt-- }()

	cond, err := a.analyzeExpr(stmt.Cond)
	if err != nil {
		return stmt, err
	}
	stmt.Cond = cond

	body, err := a.analyzeStmt(stmt.Body)
	if err != nil {
		return stmt, err
	}
	stmt.Body = body

	return stmt, nil
}

func (a *analyzer) analyzeVarDecl(decl stmts.VarDecl) (ast.Stmt, error) {
	_, ok := a.symTable.Insert(decl.Var.String(), decl.Mutable)
	if !ok {
		return nil, RedeclaredIdentError(decl.Var)
	}

	if decl.Value != nil {
		value, err := a.analyzeExpr(decl.Value)
		if err != nil {
			return nil, err
		}
		decl.Value = value
	}

	return decl, nil
}

func (a *analyzer) analyzeReturn(stmt stmts.Return) (ast.Stmt, error) {
	if a.funcCnt == 0 {
		return nil, ErrReturnOutsideFunc
	}

	if stmt.Value != nil {
		value, err := a.analyzeExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
		stmt.Value = value
	}

	return stmt, nil
}

func (a *analyzer) analyzeFuncDecl(decl stmts.FuncDecl) (ast.Stmt, error) {
	nameSym, ok := a.symTable.Insert(decl.Name.String(), false)
	if !ok {
		return nil, RedeclaredIdentError(decl.Name)
	}

	decl.Name = exprs.Ident{nameSym}

	a.funcCnt++
	defer func() { a.funcCnt-- }()

	decl.Container = a.symTable.OpenContainer()
	defer a.symTable.CloseContainer()

	for _, param := range decl.Params {
		_, ok := a.symTable.Insert(param.String(), true)
		if !ok {
			return nil, RedeclaredIdentError(param)
		}
	}

	body, err := a.analyzeStmtBlock(decl.Body)
	if err != nil {
		return nil, err
	}

	decl.Body = body

	return decl, nil
}

func (a *analyzer) analyzeIfStmt(stmt stmts.If) (ast.Stmt, error) {
	cond, err := a.analyzeExpr(stmt.Cond)
	if err != nil {
		return nil, err
	}
	stmt.Cond = cond

	cons, err := a.analyzeStmt(stmt.Cons)
	if err != nil {
		return nil, err
	}
	stmt.Cons = cons

	if stmt.Alt != nil {
		alt, err := a.analyzeStmt(stmt.Alt)
		if err != nil {
			return nil, err
		}
		stmt.Alt = alt
	}

	return stmt, nil
}

func (a *analyzer) analyzeAssignment(stmt stmts.Assignment) (ast.Stmt, error) {
	v, err := a.analyzeExpr(stmt.Var)
	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case exprs.Ident, exprs.CollIndex:
	default:
		return nil, InvalidAssignmentError{}
	}

	stmt.Var = v

	v, err = a.analyzeExpr(stmt.Value)
	if err != nil {
		return nil, err
	}

	stmt.Value = v

	return stmt, nil
}

func (a *analyzer) analyzeCompoundAssignment(stmt stmts.CompoundAssignment) (ast.Stmt, error) {
	v, err := a.analyzeExpr(stmt.Var)
	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case exprs.Ident:
	case exprs.CollIndex:
		break
	default:
		return nil, InvalidAssignmentError{}
	}

	stmt.Var = v

	v, err = a.analyzeExpr(stmt.Value)
	if err != nil {
		return nil, err
	}

	stmt.Value = v

	return stmt, nil
}
