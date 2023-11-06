package lexer

import "unicode"

func (lexer *lexer) nextChar() rune {
	for {
		ch, _, err := lexer.codeReader.ReadRune()
		if err != nil {
			return 0
		}

		return ch
	}
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
