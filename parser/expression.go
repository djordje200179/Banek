package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strconv"
)

func (parser *Parser) parseExpression(precedence OperatorPrecedence) (ast.Expression, error) {
	prefixParser, ok := parser.prefixParsers[parser.currentToken.Type]
	if !ok {
		return nil, UnknownTokenError{TokenType: parser.currentToken.Type}
	}

	leftExp, err := prefixParser()
	if err != nil {
		return nil, err
	}

	for parser.nextToken.Type != tokens.SemiColon && precedence < infixOperatorPrecedences[parser.nextToken.Type] {
		infixParser, ok := parser.infixParsers[parser.nextToken.Type]
		if !ok {
			return leftExp, nil
		}

		parser.fetchToken()

		leftExp, err = infixParser(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (parser *Parser) parseIdentifier() (ast.Expression, error) {
	return expressions.Identifier{Name: parser.currentToken.Literal}, nil
}

func (parser *Parser) parseIntegerLiteral() (ast.Expression, error) {
	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	return expressions.IntegerLiteral{Value: value}, nil
}

func (parser *Parser) parseBooleanLiteral() (ast.Expression, error) {
	value, err := strconv.ParseBool(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}

	return expressions.BooleanLiteral{Value: value}, nil
}

func (parser *Parser) parsePrefixedExpression() (ast.Expression, error) {
	expression := expressions.PrefixedExpression{
		Operator: parser.currentToken,
	}

	parser.fetchToken()

	var err error
	expression.Wrapped, err = parser.parseExpression(Prefix)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (parser *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := expressions.InfixExpression{
		Left:     left,
		Operator: parser.currentToken,
	}

	precedence, ok := infixOperatorPrecedences[parser.currentToken.Type]
	if !ok {
		return nil, UnknownTokenError{TokenType: parser.currentToken.Type}
	}

	parser.fetchToken()

	var err error
	expression.Right, err = parser.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (parser *Parser) parseGroupedExpression() (ast.Expression, error) {
	parser.fetchToken()

	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err = parser.expectNextToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	return expression, nil
}

func (parser *Parser) parseIfExpression() (ast.Expression, error) {
	var expression expressions.IfExpression
	var err error

	if err = parser.expectNextToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	expression.Condition, err = parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err = parser.expectNextToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	expression.Consequence, err = parser.parseStatement()
	if err != nil {
		return nil, err
	}

	if parser.nextToken.Type == tokens.Else {
		parser.fetchToken()
		parser.fetchToken()

		expression.Alternative, err = parser.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return expression, nil
}

func (parser *Parser) parseFunctionLiteral() (ast.Expression, error) {
	var err error

	if err = parser.expectNextToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	var literal expressions.FunctionLiteral

	literal.Parameters, err = parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	literal.Body, err = parser.parseStatement()
	if err != nil {
		return nil, err
	}

	return literal, nil
}

func (parser *Parser) parseFunctionParameters() ([]expressions.Identifier, error) {
	parser.fetchToken()

	if parser.currentToken.Type == tokens.RightParenthesis {
		parser.fetchToken()
		return nil, nil
	}

	var parameters []expressions.Identifier

	identifier, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}
	parameters = append(parameters, identifier.(expressions.Identifier))

	for parser.nextToken.Type == tokens.Comma {
		parser.fetchToken()
		parser.fetchToken()

		identifier, err = parser.parseIdentifier()
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, identifier.(expressions.Identifier))
	}

	if err = parser.expectNextToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return parameters, nil
}

func (parser *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	arguments, err := parser.parseCallArguments()
	if err != nil {
		return nil, err
	}

	return expressions.FunctionCall{Function: function, Arguments: arguments}, nil
}

func (parser *Parser) parseCallArguments() ([]ast.Expression, error) {
	parser.fetchToken()

	if parser.currentToken.Type == tokens.RightParenthesis {
		parser.fetchToken()
		return nil, nil
	}

	var arguments []ast.Expression

	argument, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, argument)

	for parser.nextToken.Type == tokens.Comma {
		parser.fetchToken()
		parser.fetchToken()

		argument, err = parser.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argument)
	}

	if err = parser.expectNextToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return arguments, nil
}
