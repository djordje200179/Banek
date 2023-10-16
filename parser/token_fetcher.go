package parser

import "banek/tokens"

func (parser *parser) fetchToken() {
	if parser.nextToken.Type == tokens.EOF {
		parser.currentToken = tokens.Token{Type: tokens.EOF}
		return
	}

	parser.currentToken, parser.nextToken = parser.nextToken, <-parser.tokenChannel
}

func (parser *parser) assertToken(tokenType tokens.TokenType) error {
	if parser.currentToken.Type != tokenType {
		return UnexpectedTokenError{Expected: tokenType, Got: parser.currentToken.Type}
	}

	return nil
}
