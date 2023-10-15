package interpreter

import (
	"banek/interpreter/results"
	"banek/parser"
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

func (interpreter *Interpreter) Eval(statementsChan <-chan parser.ParsedStatement, bufferSize int) <-chan Result {
	resultsChan := make(chan Result, bufferSize)

	go interpreter.evalThread(statementsChan, resultsChan)

	return resultsChan
}

func (interpreter *Interpreter) evalThread(statementsChan <-chan parser.ParsedStatement, resultsChan chan<- Result) {
	for statement := range statementsChan {
		if statement.Error != nil {
			resultsChan <- results.Error{Err: statement.Error}
			continue
		}

		result, err := interpreter.evalStatement(interpreter.globalEnv, statement.Statement)
		if err != nil {
			resultsChan <- results.Error{Err: err}
			continue
		}

		resultsChan <- result
	}

	close(resultsChan)
}
