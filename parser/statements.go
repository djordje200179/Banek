package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/tokens"
)

func (parser *Parser) parseStatement() (ast.Statement, error) {
	statementParser := parser.statementParsers[parser.currentToken.Type]
	if statementParser == nil {
		statementParser = parser.parseExpressionStatement
	}

	return statementParser()
}

func (parser *Parser) parseVariableDeclarationStatement() (ast.Statement, error) {
	isConst := parser.currentToken.Type == tokens.Const

	if err := parser.expectNextToken(tokens.Identifier); err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := parser.expectNextToken(tokens.Assign); err != nil {
		return nil, err
	}

	parser.fetchToken()

	value, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.VariableDeclaration{Const: isConst, Name: name.(expressions.Identifier), Value: value}, nil
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	parser.fetchToken()

	value, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.Return{Value: value}, nil
}

func (parser *Parser) parseExpressionStatement() (ast.Statement, error) {
	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.Expression{Expression: expression}, nil
}

func (parser *Parser) parseBlockStatement() (ast.Statement, error) {
	parser.fetchToken()

	var containedStatements []ast.Statement
	for parser.currentToken.Type != tokens.RightBrace {
		singleStatement, err := parser.parseStatement()
		if err != nil {
			return nil, err
		}

		containedStatements = append(containedStatements, singleStatement)

		parser.fetchToken()
	}

	return statements.Block{Statements: containedStatements}, nil
}

func (parser *Parser) parseFunctionStatement() (ast.Statement, error) {
	if err := parser.expectNextToken(tokens.Identifier); err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	body, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	return statements.Function{Name: name.(expressions.Identifier), Parameters: parameters, Body: body.(statements.Block)}, nil
}
