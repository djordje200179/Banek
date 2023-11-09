package parser

import (
	"banek/ast"
	"banek/ast/stmts"
	"banek/tokens"
	"runtime"
)

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

func Parse(tokenChan <-chan tokens.Token, bufferSize int) <-chan ast.Stmt {
	stmtChan := make(chan ast.Stmt, bufferSize)

	go parsingThread(tokenChan, stmtChan)

	return stmtChan
}

func parsingThread(tokenChan <-chan tokens.Token, stmtChan chan<- ast.Stmt) {
	runtime.LockOSThread()

	parser := &parser{tokenChan: tokenChan}
	parser.initHandlers()

	parser.fetchToken()

	for parser.currToken.Type != tokens.EOF {
		stmt, err := parser.parseStmt()
		if err != nil {
			stmt = stmts.Invalid{Err: err}
		}

		stmtChan <- stmt
	}

	close(stmtChan)
}

func (parser *parser) initHandlers() {
	parser.prefixExprHandlers = map[tokens.Type]prefixExprHandler{
		tokens.Identifier: parser.parseIdentifier,

		tokens.Integer:   parser.parseInteger,
		tokens.Boolean:   parser.parseBoolean,
		tokens.String:    parser.parseString,
		tokens.Undefined: parser.parseUndefined,

		tokens.Minus: parser.parseUnaryOp,
		tokens.Bang:  parser.parseUnaryOp,

		tokens.LeftParen: parser.parseGroupedExpr,

		tokens.If:          parser.parseIfExpr,
		tokens.VerticalBar: parser.parseFuncLiteral,

		tokens.LeftBracket: parser.parseArray,
	}

	parser.infixExprHandlers = map[tokens.Type]infixExprHandler{
		tokens.Equals:        parser.parseBinaryOp,
		tokens.NotEquals:     parser.parseBinaryOp,
		tokens.Less:          parser.parseBinaryOp,
		tokens.Greater:       parser.parseBinaryOp,
		tokens.LessEquals:    parser.parseBinaryOp,
		tokens.GreaterEquals: parser.parseBinaryOp,

		tokens.Plus:     parser.parseBinaryOp,
		tokens.Minus:    parser.parseBinaryOp,
		tokens.Asterisk: parser.parseBinaryOp,
		tokens.Slash:    parser.parseBinaryOp,
		tokens.Modulo:   parser.parseBinaryOp,
		tokens.Caret:    parser.parseBinaryOp,

		tokens.Assign:         parser.parseAssignment,
		tokens.PlusAssign:     parser.parseAssignment,
		tokens.MinusAssign:    parser.parseAssignment,
		tokens.AsteriskAssign: parser.parseAssignment,
		tokens.SlashAssign:    parser.parseAssignment,
		tokens.ModuloAssign:   parser.parseAssignment,
		tokens.CaretAssign:    parser.parseAssignment,

		tokens.LeftParen:   parser.parseFuncCall,
		tokens.LeftBracket: parser.parseIndexExpr,
	}

	parser.stmtHandlers = map[tokens.Type]stmtHandler{
		tokens.Let: parser.parseVarDeclaration,

		tokens.Return:    parser.parseReturn,
		tokens.LeftBrace: parser.parseBlock,
		tokens.Func:      parser.parseFuncStmt,
		tokens.If:        parser.parseIfStmt,
		tokens.While:     parser.parseWhile,
	}
}
