package interpreter

import (
	"banek/ast"
	"banek/interpreter/results"
	"fmt"
)

type Interpreter struct {
	globalEnv *environment
}

func New() *Interpreter {
	interpreter := &Interpreter{
		globalEnv: newEnvironment(nil),
	}

	return interpreter
}

type Result interface {
	fmt.Stringer
}

func (interpreter *Interpreter) Eval(statementsChan <-chan ast.Statement, bufferSize int) <-chan Result {
	resultsChan := make(chan Result, bufferSize)

	go interpreter.evalThread(statementsChan, resultsChan)

	return resultsChan
}

func (interpreter *Interpreter) evalThread(statementsChan <-chan ast.Statement, resultsChan chan<- Result) {
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
