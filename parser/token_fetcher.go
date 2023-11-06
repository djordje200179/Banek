package parser

import "banek/tokens"

func (parser *parser) fetchToken() {
	nextToken, ok := <-parser.tokenChan
	if !ok {
		parser.currToken = tokens.Token{Type: tokens.EOF}
		return
	}

	parser.currToken = nextToken
}

func (parser *parser) assertToken(tokenType tokens.Type) error {
	if parser.currToken.Type != tokenType {
		return ErrUnexpectedToken{Expected: tokenType, Got: parser.currToken.Type}
	}

	return nil
}
