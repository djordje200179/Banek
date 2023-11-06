package lexer

import (
	"banek/tokens"
	"strings"
	"unicode"
)

func (lexer *lexer) readIdentifier() string {
	var sb strings.Builder

	firstChar, _, _ := lexer.codeReader.ReadRune()
	sb.WriteRune(firstChar)

	for {
		ch := lexer.nextChar()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			sb.WriteRune(ch)
		} else if ch == 0 {
			break
		} else {
			_ = lexer.codeReader.UnreadRune()
			break
		}
	}

	return sb.String()
}

func (lexer *lexer) readNumber() string {
	var sb strings.Builder

	for {
		ch := lexer.nextChar()
		if unicode.IsDigit(ch) {
			sb.WriteRune(ch)
		} else if ch == 0 {
			break
		} else {
			_ = lexer.codeReader.UnreadRune()
			break
		}
	}

	return sb.String()
}

func (lexer *lexer) readString() tokens.Token {
	var sb strings.Builder

	for {
		ch := lexer.nextChar()
		if ch == '"' {
			break
		} else if ch == 0 {
			return tokens.Token{Type: tokens.Illegal}
		}

		sb.WriteRune(ch)
	}

	return tokens.Token{Type: tokens.String, Literal: sb.String()}
}
