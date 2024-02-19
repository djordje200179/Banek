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
	reader *bufio.Reader
}

func lexingThread(reader io.Reader, tokenChan chan<- tokens.Token) {
	runtime.LockOSThread()

	var bufferedReader *bufio.Reader
	var ok bool
	if bufferedReader, ok = reader.(*bufio.Reader); !ok {
		bufferedReader = bufio.NewReader(reader)
	}

	lexer := lexer{
		reader: bufferedReader,
	}

	for {
		token, err := lexer.nextToken()
		if err != nil {
			close(tokenChan)
			panic(err)
		}

		tokenChan <- token

		if token.Type == tokens.EOF {
			break
		}
	}

	close(tokenChan)
}
