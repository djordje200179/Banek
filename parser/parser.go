package parser

import (
	"banek/ast"
	"banek/tokens"
	"runtime"
)

func Parse(tokenChan <-chan tokens.Token, bufferSize int) <-chan ast.Stmt {
	stmtChan := make(chan ast.Stmt, bufferSize)

	go parsingThread(tokenChan, stmtChan)

	return stmtChan
}

type (
	prefixExprHandler func() (ast.Expr, error)
	infixExprHandler  func(ast.Expr) (ast.Expr, error)
	stmtHandler       func() (ast.Stmt, error)
)

type parser struct {
	tokenChan <-chan tokens.Token

	currToken tokens.Token

	prefixExprHandlers map[tokens.Type]prefixExprHandler
	infixExprHandlers  map[tokens.Type]infixExprHandler
	stmtHandlers       map[tokens.Type]stmtHandler
}

func parsingThread(tokenChan <-chan tokens.Token, stmtChan chan<- ast.Stmt) {
	runtime.LockOSThread()

	parser := &parser{tokenChan: tokenChan}
	parser.initHandlers()

	parser.fetchToken()

	for parser.currToken.Type != tokens.EOF {
		stmt, err := parser.parseStmt()
		if err != nil {
			close(stmtChan)
			panic(err)
		}

		stmtChan <- stmt
	}

	close(stmtChan)
}

func (p *parser) initHandlers() {
	p.prefixExprHandlers = map[tokens.Type]prefixExprHandler{
		tokens.Ident: p.parseIdent,

		tokens.Int:       p.parseInt,
		tokens.Bool:      p.parseBool,
		tokens.String:    p.parseString,
		tokens.Undefined: p.parseUndefined,

		tokens.Minus:  p.parseUnaryOp,
		tokens.Bang:   p.parseUnaryOp,
		tokens.LArrow: p.parseUnaryOp,

		tokens.LParen: p.parseGroupedExpr,

		tokens.If:   p.parseIfExpr,
		tokens.VBar: p.parseFuncLiteral,

		tokens.LBracket: p.parseArray,
	}

	p.infixExprHandlers = map[tokens.Type]infixExprHandler{
		tokens.Equals:        p.parseBinaryOp,
		tokens.NotEquals:     p.parseBinaryOp,
		tokens.Less:          p.parseBinaryOp,
		tokens.Greater:       p.parseBinaryOp,
		tokens.LessEquals:    p.parseBinaryOp,
		tokens.GreaterEquals: p.parseBinaryOp,

		tokens.Plus:     p.parseBinaryOp,
		tokens.Minus:    p.parseBinaryOp,
		tokens.Asterisk: p.parseBinaryOp,
		tokens.Slash:    p.parseBinaryOp,
		tokens.Percent:  p.parseBinaryOp,
		tokens.LArrow:   p.parseBinaryOp,

		tokens.Assign:         p.parseBinaryOp,
		tokens.PlusAssign:     p.parseBinaryOp,
		tokens.MinusAssign:    p.parseBinaryOp,
		tokens.AsteriskAssign: p.parseBinaryOp,
		tokens.SlashAssign:    p.parseBinaryOp,
		tokens.PercentAssign:  p.parseBinaryOp,

		tokens.LParen:   p.parseFuncCall,
		tokens.LBracket: p.parseIndexExpr,
	}

	p.stmtHandlers = map[tokens.Type]stmtHandler{
		tokens.Let: p.parseVarDecl,

		tokens.Return: p.parseReturn,
		tokens.LBrace: p.parseBlock,
		tokens.Func:   p.parseFuncStmt,
		tokens.If:     p.parseIfStmt,
		tokens.While:  p.parseWhile,
	}
}
