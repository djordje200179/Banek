package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/tokens"
)

func (parser *parser) parseStatement() (ast.Statement, error) {
	statementParser := parser.statementParsers[parser.currentToken.Type]
	if statementParser == nil {
		statementParser = parser.parseExpressionStatement
	}

	return statementParser()
}

func (parser *parser) parseVariableDeclarationStatement() (ast.Statement, error) {
	parser.fetchToken()

	var isMutable bool
	if parser.currentToken.Type == tokens.Mut {
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

	value, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return statements.VariableDeclaration{Mutable: isMutable, Name: name.(expressions.Identifier), Value: value}, nil
}

func (parser *parser) parseReturnStatement() (ast.Statement, error) {
	parser.fetchToken()

	value, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.SemiColon); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return statements.Return{Value: value}, nil
}

func (parser *parser) parseExpressionStatement() (ast.Statement, error) {
	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if parser.currentToken.Type == tokens.SemiColon {
		parser.fetchToken()
	}

	return statements.Expression{Expression: expression}, nil
}

func (parser *parser) parseBlockStatement() (ast.Statement, error) {
	parser.fetchToken()

	var containedStatements []ast.Statement
	for parser.currentToken.Type != tokens.RightBrace {
		singleStatement, err := parser.parseStatement()
		if err != nil {
			return nil, err
		}

		containedStatements = append(containedStatements, singleStatement)
	}

	if err := parser.assertToken(tokens.RightBrace); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return statements.Block{Statements: containedStatements}, nil
}

func (parser *parser) parseFunctionStatement() (ast.Statement, error) {
	parser.fetchToken()

	if err := parser.assertToken(tokens.Identifier); err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	body, err := parser.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return statements.Function{Name: name.(expressions.Identifier), Parameters: parameters, Body: body.(statements.Block)}, nil
}

func (parser *parser) parseIfStatement() (ast.Statement, error) {
	parser.fetchToken()

	condition, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	parser.fetchToken()

	consequence, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	var alternative ast.Statement
	if parser.currentToken.Type == tokens.Else {
		parser.fetchToken()

		alternative, err = parser.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return statements.If{Condition: condition, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *parser) parseWhileStatement() (ast.Statement, error) {
	parser.fetchToken()

	condition, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Do); err != nil {
		return nil, err
	}

	parser.fetchToken()

	body, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	return statements.While{Condition: condition, Body: body}, nil
}
