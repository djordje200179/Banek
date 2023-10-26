package lexer

import (
	"banek/tokens"
	"bufio"
	"io"
	"runtime"
)

func Tokenize(reader io.Reader, bufferSize int) <-chan tokens.Token {
	tokenChannel := make(chan tokens.Token, bufferSize)

	go lexingThread(reader, tokenChannel)

	return tokenChannel
}

type lexer struct {
	reader *bufio.Reader
}

func lexingThread(reader io.Reader, ch chan<- tokens.Token) {
	runtime.LockOSThread()

	lexer := lexer{
		reader: bufio.NewReader(reader),
	}

	for {
		token := lexer.nextToken()
		ch <- token

		if token.Type == tokens.EOF {
			close(ch)
			return
		}
	}
}
