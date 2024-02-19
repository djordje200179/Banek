package lexer

import (
	"banek/tokens"
	"strings"
	"unicode"
)

func (l *lexer) readIdent() (tokens.Token, error) {
	var sb strings.Builder

	for {
		ch, err := l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			_ = l.reader.UnreadRune()
			break
		}

		sb.WriteRune(ch)
	}

	tokenType := tokens.LookupIdent(sb.String())
	return tokens.Token{Type: tokenType, Literal: sb.String()}, nil
}

func (l *lexer) readNum() (tokens.Token, error) {
	var sb strings.Builder

	for {
		ch, err := l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		if !unicode.IsDigit(ch) {
			_ = l.reader.UnreadRune()
			break
		}

		sb.WriteRune(ch)
	}

	return tokens.Token{Type: tokens.Int, Literal: sb.String()}, nil
}

func (l *lexer) readString() (tokens.Token, error) {
	var sb strings.Builder

	_, _ = l.nextChar()

	for {
		ch, err := l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		if ch == '"' {
			break
		}

		sb.WriteRune(ch)
	}

	return tokens.Token{Type: tokens.String, Literal: sb.String()}, nil
}

func (l *lexer) readChar() (tokens.Token, error) {
	var tokenType tokens.Type

	ch, err := l.nextChar()
	if err != nil {
		return tokens.Token{}, err
	}

	switch ch {
	case '+':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.PlusAssign
		default:
			tokenType = tokens.Plus
			_ = l.reader.UnreadRune()
		}
	case '-':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '>':
			tokenType = tokens.RArrow
		case '=':
			tokenType = tokens.MinusAssign
		default:
			tokenType = tokens.Minus
			_ = l.reader.UnreadRune()
		}
	case '*':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.AsteriskAssign
		default:
			tokenType = tokens.Asterisk
			_ = l.reader.UnreadRune()
		}
	case '/':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.SlashAssign
		default:
			tokenType = tokens.Slash
			_ = l.reader.UnreadRune()
		}
	case '%':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.PercentAssign
		default:
			tokenType = tokens.Percent
			_ = l.reader.UnreadRune()
		}
	case '!':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.NotEquals
		default:
			tokenType = tokens.Bang
			_ = l.reader.UnreadRune()
		}
	case '<':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '-':
			tokenType = tokens.LArrow
		case '=':
			tokenType = tokens.LessEquals
		default:
			tokenType = tokens.Less
			_ = l.reader.UnreadRune()
		}
	case '|':
		tokenType = tokens.VBar
	case ',':
		tokenType = tokens.Comma
	case ';':
		tokenType = tokens.SemiColon
	case '(':
		tokenType = tokens.LParen
	case ')':
		tokenType = tokens.RParen
	case '{':
		tokenType = tokens.LBrace
	case '}':
		tokenType = tokens.RBrace
	case '[':
		tokenType = tokens.LBracket
	case ']':
		tokenType = tokens.RBracket
	case '=':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.Equals
		default:
			tokenType = tokens.Assign
			_ = l.reader.UnreadRune()
		}
	case '>':
		ch, err = l.nextChar()
		if err != nil {
			return tokens.Token{}, err
		}

		switch ch {
		case '=':
			tokenType = tokens.GreaterEquals
		default:
			tokenType = tokens.Greater
			_ = l.reader.UnreadRune()
		}
	default:
		return tokens.Token{Type: tokens.Illegal, Literal: string(ch)}, nil
	}

	return tokens.Token{Type: tokenType}, nil
}
