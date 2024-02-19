package parser

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/ast/stmts"
	"banek/tokens"
)

func (p *parser) parseStmt() (ast.Stmt, error) {
	for p.currToken.Type == tokens.SemiColon {
		p.fetchToken()
	}

	stmtHandler := p.stmtHandlers[p.currToken.Type]
	if stmtHandler == nil {
		stmtHandler = p.parseExprStmt
	}

	return stmtHandler()
}

func (p *parser) parseVarDecl() (ast.Stmt, error) {
	p.fetchToken()

	var isMutable bool
	if p.currToken.Type == tokens.Mut {
		isMutable = true
		p.fetchToken()
	}

	if err := p.assertToken(tokens.Ident); err != nil {
		return nil, err
	}

	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.Assign); err != nil {
		return nil, err
	}

	p.fetchToken()

	value, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	p.fetchToken()

	return stmts.VarDecl{Mutable: isMutable, Var: name.(exprs.Ident), Value: value}, nil
}

func (p *parser) parseReturn() (ast.Stmt, error) {
	p.fetchToken()

	value, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	p.fetchToken()

	return stmts.Return{Value: value}, nil
}

func (p *parser) parseExprStmt() (ast.Stmt, error) {
	expr, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if p.currToken.Type == tokens.SemiColon {
		p.fetchToken()
	}

	switch expr := expr.(type) {
	case exprs.FuncCall:
		return stmts.FuncCall(expr), nil
	case exprs.BinaryOp:
		switch expr.Operator {
		case tokens.Assign:
			return stmts.Assignment{Var: expr.Left, Value: expr.Right}, nil
		case tokens.PlusAssign, tokens.MinusAssign, tokens.AsteriskAssign, tokens.SlashAssign, tokens.PercentAssign:
			return stmts.CompoundAssignment{Var: expr.Left, Value: expr.Right, Operator: expr.Operator}, nil
		default:
			return nil, InvalidExprStmtError{expr}
		}
	}

	return nil, InvalidExprStmtError{expr}
}

func (p *parser) parseBlock() (ast.Stmt, error) {
	p.fetchToken()

	var block stmts.Block
	for p.currToken.Type != tokens.RBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		block = append(block, stmt)
	}

	if err := p.assertToken(tokens.RBrace); err != nil {
		return nil, err
	}

	p.fetchToken()

	return block, nil
}

func (p *parser) parseFuncStmt() (ast.Stmt, error) {
	p.fetchToken()

	if err := p.assertToken(tokens.Ident); err != nil {
		return nil, err
	}

	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.LParen); err != nil {
		return nil, err
	}

	p.fetchToken()

	params, err := p.parseFuncParams(tokens.RParen)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return stmts.FuncDecl{Name: name.(exprs.Ident), Params: params, Body: body.(stmts.Block)}, nil
}

func (p *parser) parseIfStmt() (ast.Stmt, error) {
	p.fetchToken()

	cond, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	p.fetchToken()

	cons, err := p.parseStmt()
	if err != nil {
		return nil, err
	}

	var alt ast.Stmt
	if p.currToken.Type == tokens.Else {
		p.fetchToken()

		alt, err = p.parseStmt()
		if err != nil {
			return nil, err
		}
	}

	return stmts.If{Cond: cond, Cons: cons, Alt: alt}, nil
}

func (p *parser) parseWhile() (ast.Stmt, error) {
	p.fetchToken()

	cond, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.Do); err != nil {
		return nil, err
	}

	p.fetchToken()

	body, err := p.parseStmt()
	if err != nil {
		return nil, err
	}

	return stmts.While{Cond: cond, Body: body}, nil
}
