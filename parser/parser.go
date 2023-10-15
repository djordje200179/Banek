package parser

import (
	"banek/ast"
	"banek/tokens"
)

type (
	prefixParser    func() (ast.Expression, error)
	infixParser     func(ast.Expression) (ast.Expression, error)
	statementParser func() (ast.Statement, error)
)

type Parser struct {
	tokenChannel <-chan tokens.Token

	currentToken, nextToken tokens.Token

	prefixParsers    map[tokens.TokenType]prefixParser
	infixParsers     map[tokens.TokenType]infixParser
	statementParsers map[tokens.TokenType]statementParser
}

func New() *Parser {
	parser := new(Parser)

	parser.prefixParsers = map[tokens.TokenType]prefixParser{
		tokens.Identifier: parser.parseIdentifier,

		tokens.Integer: parser.parseIntegerLiteral,
		tokens.Boolean: parser.parseBooleanLiteral,
		tokens.String:  parser.parseStringLiteral,

		tokens.Minus: parser.parsePrefixOperation,
		tokens.Bang:  parser.parsePrefixOperation,

		tokens.LeftParenthesis: parser.parseGroupedExpression,

		tokens.If:             parser.parseIfExpression,
		tokens.LambdaFunction: parser.parseFunctionLiteral,

		tokens.LeftBracket: parser.parseArrayLiteral,
	}

	parser.infixParsers = map[tokens.TokenType]infixParser{
		tokens.Equals:              parser.parseInfixOperation,
		tokens.NotEquals:           parser.parseInfixOperation,
		tokens.LessThan:            parser.parseInfixOperation,
		tokens.GreaterThan:         parser.parseInfixOperation,
		tokens.LessThanOrEquals:    parser.parseInfixOperation,
		tokens.GreaterThanOrEquals: parser.parseInfixOperation,
		tokens.Plus:                parser.parseInfixOperation,
		tokens.Minus:               parser.parseInfixOperation,
		tokens.Asterisk:            parser.parseInfixOperation,
		tokens.Slash:               parser.parseInfixOperation,

		tokens.Assign: parser.parseVariableAssignment,

		tokens.LeftParenthesis: parser.parseCallExpression,
		tokens.LeftBracket:     parser.parseIndexExpression,
	}

	parser.statementParsers = map[tokens.TokenType]statementParser{
		tokens.Let:   parser.parseVariableDeclarationStatement,
		tokens.Const: parser.parseVariableDeclarationStatement,

		tokens.Return:    parser.parseReturnStatement,
		tokens.LeftBrace: parser.parseBlockStatement,
		tokens.Function:  parser.parseFunctionStatement,
		tokens.If:        parser.parseIfStatement,
	}

	return parser
}

func (parser *Parser) fetchToken() {
	if parser.nextToken.Type == tokens.EOF {
		parser.currentToken = tokens.Token{Type: tokens.EOF}
		return
	}

	parser.currentToken, parser.nextToken = parser.nextToken, <-parser.tokenChannel
}

func (parser *Parser) assertToken(tokenType tokens.TokenType) error {
	if parser.currentToken.Type != tokenType {
		return UnexpectedTokenError{Expected: tokenType, Got: parser.currentToken.Type}
	}

	return nil
}

type ParsedStatement struct {
	Statement ast.Statement
	Error     error
}

func (parser *Parser) Parse(tokenChannel <-chan tokens.Token, bufferSize int) <-chan ParsedStatement {
	parser.tokenChannel = tokenChannel

	parser.fetchToken()
	parser.fetchToken()

	statementChannel := make(chan ParsedStatement, bufferSize)

	go parser.parsingThread(statementChannel)

	return statementChannel
}

func (parser *Parser) parsingThread(ch chan<- ParsedStatement) {
	for parser.currentToken.Type != tokens.EOF {
		statement, err := parser.parseStatement()
		if err != nil {
			ch <- ParsedStatement{Error: err}
			continue
		}

		ch <- ParsedStatement{Statement: statement}
	}

	close(ch)
}
