package parser

import (
	"banek/ast"
	"banek/lexer"
	"banek/tokens"
	"errors"
)

type (
	prefixParser    func() (ast.Expression, error)
	infixParser     func(ast.Expression) (ast.Expression, error)
	statementParser func() (ast.Statement, error)
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken, nextToken tokens.Token

	prefixParsers    map[tokens.TokenType]prefixParser
	infixParsers     map[tokens.TokenType]infixParser
	statementParsers map[tokens.TokenType]statementParser
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: lexer,
	}

	parser.fetchToken()
	parser.fetchToken()

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
	parser.currentToken, parser.nextToken = parser.nextToken, parser.lexer.NextToken()
}

func (parser *Parser) expectNextToken(tokenType tokens.TokenType) error {
	if parser.nextToken.Type != tokenType {
		return UnexpectedTokenError{Expected: tokenType, Got: parser.nextToken.Type}
	}

	parser.fetchToken()
	return nil
}

func (parser *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	var wrappedErrors error

	for ; parser.currentToken.Type != tokens.EOF; parser.fetchToken() {
		statement, err := parser.parseStatement()
		if err != nil {
			wrappedErrors = errors.Join(wrappedErrors, err)
			continue
		}

		program.Statements = append(program.Statements, statement)
	}

	return program, wrappedErrors
}
