package lexer

import (
	"banek/tokens"
	"bufio"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	reader *bufio.Reader
}

func New(reader io.Reader) *Lexer {
	lexer := &Lexer{
		reader: bufio.NewReader(reader),
	}

	return lexer
}

func (lexer *Lexer) Tokenize(bufferSize int) <-chan tokens.Token {
	tokenChannel := make(chan tokens.Token, bufferSize)

	go lexer.lexingThread(tokenChannel)

	return tokenChannel
}

func (lexer *Lexer) lexingThread(ch chan<- tokens.Token) {
	for {
		token := lexer.nextToken()
		ch <- token

		if token.Type == tokens.EOF {
			close(ch)
			return
		}
	}
}

func (lexer *Lexer) nextToken() tokens.Token {
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

func (lexer *Lexer) readIdentifier() string {
	var sb strings.Builder

	firstChar, _, _ := lexer.reader.ReadRune()
	sb.WriteRune(firstChar)

	for {
		ch := lexer.nextChar()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			sb.WriteRune(ch)
		} else if ch == 0 {
			break
		} else {
			_ = lexer.reader.UnreadRune()
			break
		}
	}

	return sb.String()
}

func (lexer *Lexer) readNumber() string {
	var sb strings.Builder

	for {
		ch := lexer.nextChar()
		if unicode.IsDigit(ch) {
			sb.WriteRune(ch)
		} else if ch == 0 {
			break
		} else {
			_ = lexer.reader.UnreadRune()
			break
		}
	}

	return sb.String()
}

func (lexer *Lexer) nextChar() rune {
	for {
		ch, _, err := lexer.reader.ReadRune()
		if err != nil {
			return 0
		}

		return ch
	}
}

func (lexer *Lexer) skipWhitespace() {
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
