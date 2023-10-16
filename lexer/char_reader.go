package lexer

import "unicode"

func (lexer *lexer) nextChar() rune {
	for {
		ch, _, err := lexer.reader.ReadRune()
		if err != nil {
			return 0
		}

		return ch
	}
}

func (lexer *lexer) skipWhitespace() {
	for {
		ch, _, err := lexer.reader.ReadRune()
		if err != nil {
			return
		}

		if !unicode.IsSpace(ch) {
			_ = lexer.reader.UnreadRune()
			return
		}
	}
}
