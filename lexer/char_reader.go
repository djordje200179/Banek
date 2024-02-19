package lexer

import (
	"errors"
	"io"
	"unicode"
)

func (l *lexer) nextChar() (rune, error) {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, nil
		}

		return 0, err
	}

	return ch, nil
}

func (l *lexer) skipBlank() error {
	for {
		ch, err := l.nextChar()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(ch) {
			break
		}
	}

	_ = l.reader.UnreadRune()
	return nil
}

func (l *lexer) skipLine() error {
	for {
		ch, err := l.nextChar()
		if err != nil {
			return err
		}

		if ch == '\n' || ch == '\r' {
			break
		}
	}

	return nil
}
