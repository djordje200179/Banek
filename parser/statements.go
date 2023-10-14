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
	var statement statements.VariableDeclarationStatement
	var err error

	statement.Const = parser.currentToken.Type == tokens.Const

	err = parser.expectNextToken(tokens.Identifier)
	if err != nil {
		return nil, err
	}

	statement.Name = expressions.Identifier{Name: parser.currentToken.Literal}

	err = parser.expectNextToken(tokens.Assign)
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	statement.Value, err = parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon {
		parser.fetchToken()
	}

	return statement, nil
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	var statement statements.ReturnStatement
	var err error

	parser.fetchToken()

	statement.Value, err = parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon {
		parser.fetchToken()
	}

	return statement, err
}

func (parser *Parser) parseExpressionStatement() (ast.Statement, error) {
	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if parser.nextToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.ExpressionStatement{Expression: expression}, nil
}

func (parser *Parser) parseBlockStatement() (ast.Statement, error) {
	var blockStatement statements.BlockStatement
	var err error

	parser.fetchToken()

	for parser.currentToken.Type != tokens.RightBrace {
		singleStatement, err := parser.parseStatement()
		if err != nil {
			return nil, err
		}

		blockStatement.Statements = append(blockStatement.Statements, singleStatement)

		parser.fetchToken()
	}

	return blockStatement, err
}

func (parser *Parser) parseFunctionStatement() (ast.Statement, error) {
	var function statements.FunctionStatement
	var err error

	err = parser.expectNextToken(tokens.Identifier)
	if err != nil {
		return nil, err
	}

	nameIdentifier, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	function.Name = nameIdentifier.(expressions.Identifier)
	parser.fetchToken()

	function.Parameters, err = parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	bodyStatement, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	function.Body = bodyStatement.(statements.BlockStatement)

	return function, nil
}
