package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/ast/statements"
	"banek/tokens"
	"strconv"
)

func (parser *parser) parseExpression(precedence OperatorPrecedence) (ast.Expression, error) {
	expressionParser := parser.prefixParsers[parser.currentToken.Type]
	if expressionParser == nil {
		return nil, ErrUnknownToken{TokenType: parser.currentToken.Type}
	}

	leftExpression, err := expressionParser()
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type != tokens.SemiColon && precedence < infixOperatorPrecedences[parser.currentToken.Type] {
		expressionParser := parser.infixParsers[parser.currentToken.Type]
		if expressionParser == nil {
			return leftExpression, nil
		}

		leftExpression, err = expressionParser(leftExpression)
		if err != nil {
			return nil, err
		}
	}

	return leftExpression, nil
}

func (parser *parser) parseIdentifier() (ast.Expression, error) {
	literal := parser.currentToken.Literal

	parser.fetchToken()

	return expressions.Identifier(literal), nil
}

func (parser *parser) parseIntegerLiteral() (ast.Expression, error) {
	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.IntegerLiteral(value), nil
}

func (parser *parser) parseBooleanLiteral() (ast.Expression, error) {
	value, err := strconv.ParseBool(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.BooleanLiteral(value), nil
}

func (parser *parser) parseStringLiteral() (ast.Expression, error) {
	literal := parser.currentToken.Literal

	parser.fetchToken()

	return expressions.StringLiteral(literal), nil
}

func (parser *parser) parseArrayLiteral() (ast.Expression, error) {
	parser.fetchToken()

	var elements expressions.ArrayLiteral

	if parser.currentToken.Type == tokens.RightBracket {
		parser.fetchToken()
		return elements, nil
	}

	argument, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	elements = append(elements, argument)

	for parser.currentToken.Type == tokens.Comma {
		parser.fetchToken()

		argument, err = parser.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}

		elements = append(elements, argument)
	}

	if err := parser.assertToken(tokens.RightBracket); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return elements, nil
}

func (parser *parser) parsePrefixOperation() (ast.Expression, error) {
	operator := parser.currentToken

	parser.fetchToken()

	operand, err := parser.parseExpression(Prefix)
	if err != nil {
		return nil, err
	}

	return expressions.PrefixOperation{Operator: operator, Operand: operand}, nil
}

func (parser *parser) parseInfixOperation(left ast.Expression) (ast.Expression, error) {
	operator := parser.currentToken

	precedence, ok := infixOperatorPrecedences[operator.Type]
	if !ok {
		return nil, ErrUnknownToken{TokenType: operator.Type}
	}

	parser.fetchToken()

	right, err := parser.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expressions.InfixOperation{Left: left, Operator: operator, Right: right}, nil
}

func (parser *parser) parseAssignment(variable ast.Expression) (ast.Expression, error) {
	parser.fetchToken()

	value, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	return expressions.Assignment{Variable: variable, Value: value}, nil
}

func (parser *parser) parseGroupedExpression() (ast.Expression, error) {
	parser.fetchToken()

	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expression, nil
}

func (parser *parser) parseIfExpression() (ast.Expression, error) {
	parser.fetchToken()

	condition, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	parser.fetchToken()

	consequence, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Else); err != nil {
		return nil, err
	}

	parser.fetchToken()

	alternative, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	return expressions.If{Condition: condition, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *parser) parseFunctionLiteral() (ast.Expression, error) {
	parser.fetchToken()

	if err := parser.assertToken(tokens.LeftParenthesis); err != nil {
		return nil, err
	}

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Minus); err != nil {
		return nil, err
	}

	parser.fetchToken()

	if err := parser.assertToken(tokens.GreaterThan); err != nil {
		return nil, err
	}

	parser.fetchToken()

	expression, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	return expressions.FunctionLiteral{Parameters: parameters, Body: statements.Expression{Expression: expression}}, nil
}

func (parser *parser) parseFunctionParameters() ([]expressions.Identifier, error) {
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

	for parser.currentToken.Type == tokens.Comma {
		parser.fetchToken()

		identifier, err = parser.parseIdentifier()
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, identifier.(expressions.Identifier))
	}

	if err := parser.assertToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return parameters, nil
}

func (parser *parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	arguments, err := parser.parseCallArguments()
	if err != nil {
		return nil, err
	}

	return expressions.FunctionCall{Function: function, Arguments: arguments}, nil
}

func (parser *parser) parseIndexExpression(collection ast.Expression) (ast.Expression, error) {
	parser.fetchToken()

	index, err := parser.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.RightBracket); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.CollectionAccess{Collection: collection, Key: index}, nil
}

func (parser *parser) parseCallArguments() ([]ast.Expression, error) {
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

	for parser.currentToken.Type == tokens.Comma {
		parser.fetchToken()

		argument, err = parser.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argument)
	}

	if err := parser.assertToken(tokens.RightParenthesis); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return arguments, nil
}
