package lexer

import (
	"banek/tokens"
	"strings"
	"unicode"
)

func (lexer *lexer) nextToken() tokens.Token {
	lexer.skipWhitespace()

	nextChar := lexer.nextChar()
	if nextChar == 0 {
		return tokens.Token{Type: tokens.EOF}
	}

	switch {
	case unicode.IsLetter(nextChar):
		_ = lexer.reader.UnreadRune()

		identifier := lexer.readIdentifier()
		tokenType := tokens.LookupIdentifier(identifier)

		return tokens.Token{Type: tokenType, Literal: identifier}
	case unicode.IsDigit(nextChar):
		_ = lexer.reader.UnreadRune()

		number := lexer.readNumber()
		return tokens.Token{Type: tokens.Integer, Literal: number}
	case nextChar == '"':
		return lexer.readString()
	}

	var possibleCharTokens []string
	for token := range tokens.CharTokens {
		if strings.HasPrefix(token, string(nextChar)) {
			possibleCharTokens = append(possibleCharTokens, token)
		}
	}

	var currentToken strings.Builder
	for {
		newToken := currentToken.String() + string(nextChar)
		var nextPossibleCharTokens []string
		for _, possibleCharToken := range possibleCharTokens {
			if strings.HasPrefix(possibleCharToken, newToken) {
				nextPossibleCharTokens = append(nextPossibleCharTokens, possibleCharToken)
			}
		}

		if len(nextPossibleCharTokens) == 0 {
			_ = lexer.reader.UnreadRune()
			break
		}

		currentToken.WriteRune(nextChar)
		nextChar = lexer.nextChar()
		if nextChar == 0 {
			continue
		}
	}

	return tokens.Token{Type: tokens.CharTokens[currentToken.String()]}
}
