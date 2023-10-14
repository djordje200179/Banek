package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/tokens"
	"strconv"
)

func (parser *Parser) parseExpression(precedence OperatorPrecedence) (ast.Expression, error) {
	expressionParser, ok := parser.prefixParsers[parser.currentToken.Type]
	if !ok {
		return nil, UnknownTokenError{TokenType: parser.currentToken.Type}
	}

	leftExpression, err := expressionParser()
	if err != nil {
		return nil, err
	}

	for parser.nextToken.Type != tokens.SemiColon && precedence < infixOperatorPrecedences[parser.nextToken.Type] {
		expressionParser, ok := parser.infixParsers[parser.nextToken.Type]
		if !ok {
			return leftExpression, nil
		}

		parser.fetchToken()

		leftExpression, err = expressionParser(leftExpression)
		if err != nil {
			return nil, err
		}
	}

	return leftExpression, nil
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

func (parser *Parser) parsePrefixOperation() (ast.Expression, error) {
	operator := parser.currentToken

	parser.fetchToken()

	operand, err := parser.parseExpression(Prefix)
	if err != nil {
		return nil, err
	}

	return expressions.PrefixOperation{Operator: operator, Operand: operand}, nil
}

func (parser *Parser) parseInfixOperation(left ast.Expression) (ast.Expression, error) {
	operator := parser.currentToken

	precedence, ok := infixOperatorPrecedences[parser.currentToken.Type]
	if !ok {
		return nil, UnknownTokenError{TokenType: parser.currentToken.Type}
	}

	parser.fetchToken()

	right, err := parser.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expressions.InfixOperation{Left: left, Operator: operator, Right: right}, nil
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

	parser.fetchToken()

	return expression, nil
}

func (parser *Parser) parseIfExpression() (ast.Expression, error) {
	if err := parser.expectNextToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	condition, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.expectNextToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	consequence, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	var alternative ast.Statement
	if parser.nextToken.Type == tokens.Else {
		parser.fetchToken()
		parser.fetchToken()

		alternative, err = parser.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return expressions.IfExpression{Condition: condition, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *Parser) parseFunctionLiteral() (ast.Expression, error) {
	if err := parser.expectNextToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	body, err := parser.parseStatement()
	if err != nil {
		return nil, err
	}

	return expressions.FunctionLiteral{Parameters: parameters, Body: body}, nil
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

	if err := parser.expectNextToken(tokens.RightParenthesis); err != nil {
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
