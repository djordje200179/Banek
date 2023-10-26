package interpreter

import (
	"banek/ast"
	"banek/exec/environments"
	"banek/interpreter/results"
	"runtime"
)

func Interpret(statementsChan <-chan ast.Statement, bufferSize int) <-chan results.Result {
	resultsChan := make(chan results.Result, bufferSize)

	go evalThread(statementsChan, resultsChan)

	return resultsChan
}

type interpreter struct {
	globalEnv environments.Environment
}

var EnvFactory environments.EnvironmentFactory = environments.NewArrayEnvironment

func evalThread(statementsChan <-chan ast.Statement, resultsChan chan<- results.Result) {
	runtime.LockOSThread()

	interpreter := &interpreter{
		globalEnv: EnvFactory(nil, 0),
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
