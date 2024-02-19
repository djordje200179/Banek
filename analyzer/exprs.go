package analyzer

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/tokens"
)

func (a *analyzer) analyzeExpr(expr ast.Expr) (ast.Expr, error) {
	switch expr := expr.(type) {
	case exprs.ArrayLiteral:
		return a.analyzeArrayLiteral(expr)
	case exprs.BinaryOp:
		return a.analyzeBinaryOp(expr)
	case exprs.BoolLiteral:
		return expr, nil
	case exprs.CollIndex:
		return a.analyzeCollIndex(expr)
	case exprs.FuncCall:
		return a.analyzeFuncCall(expr)
	case exprs.FuncLiteral:
		return a.analyzeFuncLiteral(expr)
	case exprs.Ident:
		return a.analyzeIdent(expr)
	case exprs.If:
		return a.analyzeIfExpr(expr)
	case exprs.IntLiteral:
		return expr, nil
	case exprs.StringLiteral:
		return expr, nil
	case exprs.UnaryOp:
		return a.analyzeUnaryOp(expr)
	case exprs.UndefinedLiteral:
		return expr, nil
	default:
		panic("unreachable")
	}
}

func (a *analyzer) analyzeIdent(ident exprs.Ident) (exprs.Ident, error) {
	var ok bool
	ident.Symbol, ok = a.symTable.Lookup(ident.String())
	if !ok {
		return exprs.Ident{}, UndefinedIdentError(ident)
	}

	return ident, nil
}

func (a *analyzer) analyzeArrayLiteral(arr exprs.ArrayLiteral) (ast.Expr, error) {
	var err error

	for i, elem := range arr {
		arr[i], err = a.analyzeExpr(elem)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func (a *analyzer) analyzeFuncCall(call exprs.FuncCall) (exprs.FuncCall, error) {
	var err error

	call.Func, err = a.analyzeExpr(call.Func)
	if err != nil {
		return exprs.FuncCall{}, err
	}

	for i, arg := range call.Args {
		call.Args[i], err = a.analyzeExpr(arg)
		if err != nil {
			return exprs.FuncCall{}, err
		}
	}

	return call, nil
}

func (a *analyzer) analyzeBinaryOp(op exprs.BinaryOp) (ast.Expr, error) {
	var err error

	op.Left, err = a.analyzeExpr(op.Left)
	if err != nil {
		return nil, err
	}

	op.Right, err = a.analyzeExpr(op.Right)
	if err != nil {
		return nil, err
	}

	switch op.Operator {
	case tokens.Plus, tokens.Minus, tokens.Asterisk, tokens.Slash,
		tokens.Percent, tokens.Equals, tokens.NotEquals, tokens.Less,
		tokens.Greater, tokens.LessEquals, tokens.GreaterEquals, tokens.LArrow:
	default:
		return nil, InvalidOpError(op.Operator)
	}

	return op, nil
}

func (a *analyzer) analyzeUnaryOp(op exprs.UnaryOp) (ast.Expr, error) {
	var err error

	op.Operand, err = a.analyzeExpr(op.Operand)
	if err != nil {
		return nil, err
	}

	switch op.Operator {
	case tokens.Minus, tokens.Bang, tokens.LArrow:
	default:
		return nil, InvalidOpError(op.Operator)
	}

	return op, nil
}

func (a *analyzer) analyzeCollIndex(expr exprs.CollIndex) (ast.Expr, error) {
	var err error

	expr.Coll, err = a.analyzeExpr(expr.Coll)
	if err != nil {
		return nil, err
	}

	expr.Key, err = a.analyzeExpr(expr.Key)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (a *analyzer) analyzeIfExpr(expr exprs.If) (ast.Expr, error) {
	var err error

	expr.Cond, err = a.analyzeExpr(expr.Cond)
	if err != nil {
		return nil, err
	}

	expr.Cons, err = a.analyzeExpr(expr.Cons)
	if err != nil {
		return nil, err
	}

	expr.Alt, err = a.analyzeExpr(expr.Alt)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (a *analyzer) analyzeFuncLiteral(expr exprs.FuncLiteral) (ast.Expr, error) {
	//TODO: fix

	return expr, nil
}
