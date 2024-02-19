package lexer

import (
	"banek/tokens"
	"unicode"
)

func (l *lexer) nextToken() (tokens.Token, error) {
	err := l.skipBlank()
	if err != nil {
		return tokens.Token{}, err
	}

	var nextChar rune
	for {
		nextChar, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		if nextChar != '#' {
			break
		}

		err = l.skipLine()
		if err != nil {
			return tokens.Token{}, err
		}
	}

	_ = l.reader.UnreadRune()

	switch {
	case unicode.IsLetter(nextChar):
		return l.readIdent()
	case unicode.IsDigit(nextChar):
		return l.readNum()
	case nextChar == '"':
		return l.readString()
	case nextChar == 0:
		return tokens.Token{Type: tokens.EOF}, nil
	default:
		return l.readChar()
	}
}
