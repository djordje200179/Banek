package parser

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/tokens"
)

type (
	prefixParser    func() (ast.Expression, error)
	infixParser     func(ast.Expression) (ast.Expression, error)
	statementParser func() (ast.Statement, error)
)

type parser struct {
	tokenChannel <-chan tokens.Token

	currentToken, nextToken tokens.Token

	prefixParsers    map[tokens.TokenType]prefixParser
	infixParsers     map[tokens.TokenType]infixParser
	statementParsers map[tokens.TokenType]statementParser
}

func Parse(tokenChannel <-chan tokens.Token, bufferSize int) <-chan ast.Statement {
	statementChannel := make(chan ast.Statement, bufferSize)

	go parsingThread(tokenChannel, statementChannel)

	return statementChannel
}

func parsingThread(tokenChannel <-chan tokens.Token, statementChannel chan<- ast.Statement) {
	parser := &parser{tokenChannel: tokenChannel}
	parser.initSubParsers()

	parser.fetchToken()
	parser.fetchToken()

	for parser.currentToken.Type != tokens.EOF {
		statement, err := parser.parseStatement()
		if err != nil {
			statement = statements.Error{Err: err}
		}

		statementChannel <- statement
	}

	close(statementChannel)
}

func (parser *parser) initSubParsers() {
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
		tokens.Assign:              parser.parseInfixOperation,

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
}
