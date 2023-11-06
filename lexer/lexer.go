package lexer

import (
	"banek/tokens"
	"bufio"
	"io"
	"runtime"
)

func Tokenize(codeReader io.Reader, bufferSize int) <-chan tokens.Token {
	tokenChan := make(chan tokens.Token, bufferSize)

	go lexingThread(codeReader, tokenChan)

	return tokenChan
}

type lexer struct {
	codeReader *bufio.Reader
}

func lexingThread(codeReader io.Reader, tokenChan chan<- tokens.Token) {
	runtime.LockOSThread()

	var bufferedReader *bufio.Reader
	if codeReaderBuffered, ok := codeReader.(*bufio.Reader); ok {
		bufferedReader = codeReaderBuffered
	} else {
		bufferedReader = bufio.NewReader(codeReader)
	}

	lexer := lexer{
		codeReader: bufferedReader,
	}

	for {
		token := lexer.nextToken()
		tokenChan <- token

		if token.Type == tokens.EOF {
			close(tokenChan)
			break
		}
	}
}
