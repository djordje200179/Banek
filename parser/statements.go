package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/tokens"
)

func (parser *parser) parseStmt() (ast.Statement, error) {
	for parser.currToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	stmtHandler := parser.stmtHandlers[parser.currToken.Type]
	if stmtHandler == nil {
		stmtHandler = parser.parseExprStmt
	}

	return stmtHandler()
}

func (parser *parser) parseVarDeclaration() (ast.Statement, error) {
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

	return statements.VarDecl{Mutable: isMutable, Name: name.(expressions.Identifier), Value: value}, nil
}

func (parser *parser) parseReturn() (ast.Statement, error) {
	parser.fetchToken()

	value, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return statements.Return{Value: value}, nil
}

func (parser *parser) parseExprStmt() (ast.Statement, error) {
	expression, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if parser.currToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.Expr{Expr: expression}, nil
}

func (parser *parser) parseBlock() (ast.Statement, error) {
	parser.fetchToken()

	var stmts []ast.Statement
	for parser.currToken.Type != tokens.RightBrace {
		stmt, err := parser.parseStmt()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	if err := parser.assertToken(tokens.RightBrace); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return statements.Block{Stmts: stmts}, nil
}

func (parser *parser) parseFuncStmt() (ast.Statement, error) {
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

	return statements.Func{Name: name.(expressions.Identifier), Params: params, Body: body.(statements.Block)}, nil
}

func (parser *parser) parseIfStmt() (ast.Statement, error) {
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

	var alternative ast.Statement
	if parser.currToken.Type == tokens.Else {
		parser.fetchToken()

		alternative, err = parser.parseStmt()
		if err != nil {
			return nil, err
		}
	}

	return statements.If{Cond: cond, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *parser) parseWhile() (ast.Statement, error) {
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

	return statements.While{Cond: cond, Body: body}, nil
}
