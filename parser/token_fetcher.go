package parser

import "banek/tokens"

func (parser *parser) fetchToken() {
	nextToken, ok := <-parser.tokenChannel
	if !ok {
		parser.currentToken = tokens.Token{Type: tokens.EOF}
		return
	}

	parser.currentToken = nextToken
}

func (parser *parser) assertToken(tokenType tokens.TokenType) error {
	if parser.currentToken.Type != tokenType {
		return ErrUnexpectedToken{Expected: tokenType, Got: parser.currentToken.Type}
	}

	return nil
}
