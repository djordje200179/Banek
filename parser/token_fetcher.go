package parser

import "banek/tokens"

func (p *parser) fetchToken() {
	nextToken, ok := <-p.tokenChan
	if !ok {
		return
	}

	p.currToken = nextToken
}

func (p *parser) assertToken(t tokens.Type) error {
	if p.currToken.Type != t {
		return UnexpectedTokenError{Expected: t, Got: p.currToken.Type}
	}

	return nil
}
