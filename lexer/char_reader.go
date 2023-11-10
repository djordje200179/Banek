package lexer

import "unicode"

func (lexer *lexer) nextChar() rune {
	ch, _, err := lexer.codeReader.ReadRune()
	if err != nil {
		return 0
	}

	return ch
}

func (lexer *lexer) skipBlank() {
	for {
		ch, _, err := lexer.codeReader.ReadRune()
		if err != nil {
			return
		}

		if !unicode.IsSpace(ch) {
			_ = lexer.codeReader.UnreadRune()
			return
		}
	}
}

func (lexer *lexer) skipLineComment() {
	for {
		ch, _, err := lexer.codeReader.ReadRune()
		if err != nil {
			return
		}

		if ch == '\n' || ch == 0 {
			return
		}
	}
}
