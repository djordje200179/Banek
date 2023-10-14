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
		tokens.Integer:    parser.parseIntegerLiteral,
		tokens.Boolean:    parser.parseBooleanLiteral,

		tokens.Minus: parser.parsePrefixedExpression,
		tokens.Bang:  parser.parsePrefixedExpression,

		tokens.LeftParenthesis: parser.parseGroupedExpression,

		tokens.If:             parser.parseIfExpression,
		tokens.LambdaFunction: parser.parseFunctionLiteral,
	}

	parser.infixParsers = map[tokens.TokenType]infixParser{
		tokens.Equals:              parser.parseInfixExpression,
		tokens.NotEquals:           parser.parseInfixExpression,
		tokens.LessThan:            parser.parseInfixExpression,
		tokens.GreaterThan:         parser.parseInfixExpression,
		tokens.LessThanOrEquals:    parser.parseInfixExpression,
		tokens.GreaterThanOrEquals: parser.parseInfixExpression,
		tokens.Plus:                parser.parseInfixExpression,
		tokens.Minus:               parser.parseInfixExpression,
		tokens.Asterisk:            parser.parseInfixExpression,
		tokens.Slash:               parser.parseInfixExpression,

		tokens.LeftParenthesis: parser.parseCallExpression,
	}

	parser.statementParsers = map[tokens.TokenType]statementParser{
		tokens.Var:   parser.parseVariableDeclarationStatement,
		tokens.Const: parser.parseVariableDeclarationStatement,

		tokens.Return:    parser.parseReturnStatement,
		tokens.LeftBrace: parser.parseBlockStatement,
		tokens.Function:  parser.parseFunctionStatement,
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

func (parser *Parser) expectNextToken(tokenType tokens.TokenType) error {
	if parser.nextToken.Type != tokenType {
		return UnexpectedTokenError{Expected: tokenType, Got: parser.nextToken.Type}
	}

	parser.fetchToken()
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
	for ; parser.currentToken.Type != tokens.EOF; parser.fetchToken() {
		statement, err := parser.parseStatement()
		if err != nil {
			ch <- ParsedStatement{Error: err}
			continue
		}

		ch <- ParsedStatement{Statement: statement}
	}

	close(ch)
}
