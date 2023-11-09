package parser

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/ast/stmts"
	"banek/tokens"
)

func (parser *parser) parseStmt() (ast.Stmt, error) {
	for parser.currToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	stmtHandler := parser.stmtHandlers[parser.currToken.Type]
	if stmtHandler == nil {
		stmtHandler = parser.parseExprStmt
	}

	return stmtHandler()
}

func (parser *parser) parseVarDeclaration() (ast.Stmt, error) {
	parser.fetchToken()

	var isMutable bool
	if parser.currToken.Type == tokens.Mut {
		isMutable = true
		parser.fetchToken()
	}

	if err := parser.assertToken(tokens.Identifier); err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Assign); err != nil {
		return nil, err
	}

	parser.fetchToken()

	value, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return stmts.VarDecl{Mutable: isMutable, Name: name.(exprs.Identifier), Value: value}, nil
}

func (parser *parser) parseReturn() (ast.Stmt, error) {
	parser.fetchToken()

	value, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return stmts.Return{Value: value}, nil
}

func (parser *parser) parseExprStmt() (ast.Stmt, error) {
	expression, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if parser.currToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	return stmts.Expr{Expr: expression}, nil
}

func (parser *parser) parseBlock() (ast.Stmt, error) {
	parser.fetchToken()

	var stmtsArray []ast.Stmt
	for parser.currToken.Type != tokens.RightBrace {
		stmt, err := parser.parseStmt()
		if err != nil {
			return nil, err
		}

		stmtsArray = append(stmtsArray, stmt)
	}

	if err := parser.assertToken(tokens.RightBrace); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return stmts.Block{Stmts: stmtsArray}, nil
}

func (parser *parser) parseFuncStmt() (ast.Stmt, error) {
	parser.fetchToken()

	if err := parser.assertToken(tokens.Identifier); err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.LeftParen); err != nil {
		return nil, err
	}

	parser.fetchToken()

	params, err := parser.parseFuncParams(tokens.RightParen)
	if err != nil {
		return nil, err
	}

	body, err := parser.parseBlock()
	if err != nil {
		return nil, err
	}

	return stmts.Func{Name: name.(exprs.Identifier), Params: params, Body: body.(stmts.Block)}, nil
}

func (parser *parser) parseIfStmt() (ast.Stmt, error) {
	parser.fetchToken()

	cond, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	parser.fetchToken()

	consequence, err := parser.parseStmt()
	if err != nil {
		return nil, err
	}

	var alternative ast.Stmt
	if parser.currToken.Type == tokens.Else {
		parser.fetchToken()

		alternative, err = parser.parseStmt()
		if err != nil {
			return nil, err
		}
	}

	return stmts.If{Cond: cond, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *parser) parseWhile() (ast.Stmt, error) {
	parser.fetchToken()

	cond, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Do); err != nil {
		return nil, err
	}

	parser.fetchToken()

	body, err := parser.parseStmt()
	if err != nil {
		return nil, err
	}

	return stmts.While{Cond: cond, Body: body}, nil
}
