package interpreter

import (
	"banek/ast"
	"banek/interpreter/results"
)

func Interpret(statementsChan <-chan ast.Statement, bufferSize int) <-chan results.Result {
	resultsChan := make(chan results.Result, bufferSize)

	go evalThread(statementsChan, resultsChan)

	return resultsChan
}

type interpreter struct {
	globalEnv *environment
}

func evalThread(statementsChan <-chan ast.Statement, resultsChan chan<- results.Result) {
	interpreter := &interpreter{
		globalEnv: newEnvironment(nil, 0),
	}

	for statement := range statementsChan {
		result, err := interpreter.evalStatement(interpreter.globalEnv, statement)
		if err != nil {
			resultsChan <- results.Error{Err: err}
			continue
		}

		resultsChan <- result
	}

	close(resultsChan)
}
