package parser

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/exec/objects"
	"banek/tokens"
	"strconv"
)

func (parser *parser) parseExpr(precedence OperatorPrecedence) (ast.Expression, error) {
	exprHandler := parser.prefixExprHandlers[parser.currToken.Type]
	if exprHandler == nil {
		return nil, ErrUnknownToken{TokenType: parser.currToken.Type}
	}

	leftExpr, err := exprHandler()
	if err != nil {
		return nil, err
	}

	for parser.currToken.Type != tokens.SemiColon && precedence < infixOperatorPrecedences[parser.currToken.Type] {
		exprHandler := parser.infixExprHandlers[parser.currToken.Type]
		if exprHandler == nil {
			return leftExpr, nil
		}

		leftExpr, err = exprHandler(leftExpr)
		if err != nil {
			return nil, err
		}
	}

	return leftExpr, nil
}

func (parser *parser) parseIdentifier() (ast.Expression, error) {
	literal := parser.currToken.Literal

	parser.fetchToken()

	return expressions.Identifier(literal), nil
}

func (parser *parser) parseInteger() (ast.Expression, error) {
	value, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.ConstLiteral{Value: objects.Integer(value)}, nil
}

func (parser *parser) parseBoolean() (ast.Expression, error) {
	value, err := strconv.ParseBool(parser.currToken.Literal)
	if err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.ConstLiteral{Value: objects.Boolean(value)}, nil
}

func (parser *parser) parseString() (ast.Expression, error) {
	value := parser.currToken.Literal

	parser.fetchToken()

	return expressions.ConstLiteral{Value: objects.String(value)}, nil
}

func (parser *parser) parseUndefined() (ast.Expression, error) {
	parser.fetchToken()

	return expressions.ConstLiteral{Value: objects.Undefined{}}, nil
}

func (parser *parser) parseArray() (ast.Expression, error) {
	parser.fetchToken()

	var elems expressions.ArrayLiteral

	if parser.currToken.Type == tokens.RightBracket {
		parser.fetchToken()
		return elems, nil
	}

	elem, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	elems = append(elems, elem)

	for parser.currToken.Type == tokens.Comma {
		parser.fetchToken()

		elem, err = parser.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	if err := parser.assertToken(tokens.RightBracket); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return elems, nil
}

func (parser *parser) parseUnaryOp() (ast.Expression, error) {
	opToken := parser.currToken

	parser.fetchToken()

	operand, err := parser.parseExpr(Prefix)
	if err != nil {
		return nil, err
	}

	return expressions.UnaryOp{Operation: unaryOps[opToken.Type], Operand: operand}, nil
}

func (parser *parser) parseBinaryOp(left ast.Expression) (ast.Expression, error) {
	opToken := parser.currToken

	precedence, ok := infixOperatorPrecedences[opToken.Type]
	if !ok {
		return nil, ErrUnknownToken{TokenType: opToken.Type}
	}

	parser.fetchToken()

	right, err := parser.parseExpr(precedence)
	if err != nil {
		return nil, err
	}

	return expressions.BinaryOp{Left: left, Operator: binaryOps[opToken.Type], Right: right}, nil
}

func (parser *parser) parseAssignment(variable ast.Expression) (ast.Expression, error) {
	var valueWrapper expressions.BinaryOp
	hasWrapper := false
	if parser.currToken.Type != tokens.Assign {
		op := tokens.CharTokens[parser.currToken.Type.String()[0:1]]
		valueWrapper = expressions.BinaryOp{Left: variable, Operator: binaryOps[op]}
		hasWrapper = true
	}

	parser.fetchToken()

	value, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if hasWrapper {
		valueWrapper.Right = value
		value = valueWrapper
	}

	return expressions.Assignment{Var: variable, Value: value}, nil
}

func (parser *parser) parseGroupedExpr() (ast.Expression, error) {
	parser.fetchToken()

	expr, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.RightParen); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expr, nil
}

func (parser *parser) parseIfExpr() (ast.Expression, error) {
	parser.fetchToken()

	condition, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	parser.fetchToken()

	consequence, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Else); err != nil {
		return nil, err
	}

	parser.fetchToken()

	alternative, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	return expressions.If{Cond: condition, Consequence: consequence, Alternative: alternative}, nil
}

func (parser *parser) parseFuncLiteral() (ast.Expression, error) {
	parser.fetchToken()

	params, err := parser.parseFuncParams(tokens.VerticalBar)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.Arrow); err != nil {
		return nil, err
	}

	parser.fetchToken()

	expr, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	return expressions.FuncLiteral{Params: params, Body: expr}, nil
}

func (parser *parser) parseFuncParams(end tokens.Type) ([]expressions.Identifier, error) {
	if parser.currToken.Type == end {
		parser.fetchToken()
		return nil, nil
	}

	var params []expressions.Identifier

	identifier, err := parser.parseIdentifier()
	if err != nil {
		return nil, err
	}

	params = append(params, identifier.(expressions.Identifier))

	for parser.currToken.Type == tokens.Comma {
		parser.fetchToken()

		identifier, err = parser.parseIdentifier()
		if err != nil {
			return nil, err
		}

		params = append(params, identifier.(expressions.Identifier))
	}

	if err := parser.assertToken(end); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return params, nil
}

func (parser *parser) parseFuncCall(function ast.Expression) (ast.Expression, error) {
	arguments, err := parser.parseFuncCallArgs()
	if err != nil {
		return nil, err
	}

	return expressions.FuncCall{Func: function, Args: arguments}, nil
}

func (parser *parser) parseIndexExpr(collection ast.Expression) (ast.Expression, error) {
	parser.fetchToken()

	index, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := parser.assertToken(tokens.RightBracket); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return expressions.CollIndex{Coll: collection, Key: index}, nil
}

func (parser *parser) parseFuncCallArgs() ([]ast.Expression, error) {
	parser.fetchToken()

	if parser.currToken.Type == tokens.RightParen {
		parser.fetchToken()
		return nil, nil
	}

	var args []ast.Expression

	arg, err := parser.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	args = append(args, arg)

	for parser.currToken.Type == tokens.Comma {
		parser.fetchToken()

		arg, err = parser.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	if err := parser.assertToken(tokens.RightParen); err != nil {
		return nil, err
	}

	parser.fetchToken()

	return args, nil
}
