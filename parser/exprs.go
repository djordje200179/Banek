package parser

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/symtable/symbols"
	"banek/tokens"
	"strconv"
)

func (p *parser) parseExpr(precedence OperatorPrecedence) (ast.Expr, error) {
	exprHandler := p.prefixExprHandlers[p.currToken.Type]
	if exprHandler == nil {
		return nil, InvalidOperatorError(p.currToken.Type)
	}

	leftExpr, err := exprHandler()
	if err != nil {
		return nil, err
	}

	for p.currToken.Type != tokens.SemiColon && precedence < infixOperatorPrecedences[p.currToken.Type] {
		exprHandler := p.infixExprHandlers[p.currToken.Type]
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

func (p *parser) parseIdent() (ast.Expr, error) {
	literal := p.currToken.Literal

	p.fetchToken()

	return exprs.Ident{Symbol: symbols.Ident(literal)}, nil
}

func (p *parser) parseInt() (ast.Expr, error) {
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	p.fetchToken()

	return exprs.IntLiteral(value), nil
}

func (p *parser) parseBool() (ast.Expr, error) {
	value, err := strconv.ParseBool(p.currToken.Literal)
	if err != nil {
		return nil, err
	}

	p.fetchToken()

	return exprs.BoolLiteral(value), nil
}

func (p *parser) parseString() (ast.Expr, error) {
	value := p.currToken.Literal

	p.fetchToken()

	return exprs.StringLiteral(value), nil
}

func (p *parser) parseUndefined() (ast.Expr, error) {
	p.fetchToken()

	return exprs.UndefinedLiteral{}, nil
}

func (p *parser) parseArray() (ast.Expr, error) {
	p.fetchToken()

	var elems exprs.ArrayLiteral

	if p.currToken.Type == tokens.RBracket {
		p.fetchToken()
		return elems, nil
	}

	elem, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	elems = append(elems, elem)

	for p.currToken.Type == tokens.Comma {
		p.fetchToken()

		elem, err = p.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	if err := p.assertToken(tokens.RBracket); err != nil {
		return nil, err
	}

	p.fetchToken()

	return elems, nil
}

func (p *parser) parseUnaryOp() (ast.Expr, error) {
	operator := p.currToken.Type
	p.fetchToken()

	operand, err := p.parseExpr(Prefix)
	if err != nil {
		return nil, err
	}

	return exprs.UnaryOp{Operator: operator, Operand: operand}, nil
}

func (p *parser) parseBinaryOp(left ast.Expr) (ast.Expr, error) {
	operator := p.currToken.Type
	p.fetchToken()

	precedence, ok := infixOperatorPrecedences[operator]
	if !ok {
		return nil, InvalidOperatorError(operator)
	}

	right, err := p.parseExpr(precedence)
	if err != nil {
		return nil, err
	}

	return exprs.BinaryOp{Left: left, Operator: operator, Right: right}, nil
}

func (p *parser) parseGroupedExpr() (ast.Expr, error) {
	p.fetchToken()

	expr, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.RParen); err != nil {
		return nil, err
	}

	p.fetchToken()

	return expr, nil
}

func (p *parser) parseIfExpr() (ast.Expr, error) {
	p.fetchToken()

	condition, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.Then); err != nil {
		return nil, err
	}

	p.fetchToken()

	cons, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.Else); err != nil {
		return nil, err
	}

	p.fetchToken()

	alt, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	return exprs.If{Cond: condition, Cons: cons, Alt: alt}, nil
}

func (p *parser) parseFuncLiteral() (ast.Expr, error) {
	p.fetchToken()

	params, err := p.parseFuncParams(tokens.VBar)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.RArrow); err != nil {
		return nil, err
	}

	p.fetchToken()

	expr, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	return exprs.FuncLiteral{Params: params, Body: expr}, nil
}

func (p *parser) parseFuncParams(end tokens.Type) ([]exprs.Ident, error) {
	if p.currToken.Type == end {
		p.fetchToken()
		return nil, nil
	}

	var params []exprs.Ident

	identifier, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	params = append(params, identifier.(exprs.Ident))

	for p.currToken.Type == tokens.Comma {
		p.fetchToken()

		identifier, err = p.parseIdent()
		if err != nil {
			return nil, err
		}

		params = append(params, identifier.(exprs.Ident))
	}

	if err := p.assertToken(end); err != nil {
		return nil, err
	}

	p.fetchToken()

	return params, nil
}

func (p *parser) parseFuncCall(function ast.Expr) (ast.Expr, error) {
	arguments, err := p.parseFuncCallArgs()
	if err != nil {
		return nil, err
	}

	return exprs.FuncCall{Func: function, Args: arguments}, nil
}

func (p *parser) parseIndexExpr(collection ast.Expr) (ast.Expr, error) {
	p.fetchToken()

	index, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if err := p.assertToken(tokens.RBracket); err != nil {
		return nil, err
	}

	p.fetchToken()

	return exprs.CollIndex{Coll: collection, Key: index}, nil
}

func (p *parser) parseFuncCallArgs() ([]ast.Expr, error) {
	p.fetchToken()

	if p.currToken.Type == tokens.RParen {
		p.fetchToken()
		return nil, nil
	}

	var args []ast.Expr

	arg, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	args = append(args, arg)

	for p.currToken.Type == tokens.Comma {
		p.fetchToken()

		arg, err = p.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	if err := p.assertToken(tokens.RParen); err != nil {
		return nil, err
	}

	p.fetchToken()

	return args, nil
}
