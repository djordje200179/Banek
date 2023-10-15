package evaluator

import (
	"banek/parser"
	"fmt"
)

type Evaluator struct {
	globalEnv *environment
}

func New() *Evaluator {
	evaluator := &Evaluator{
		globalEnv: newEnvironment(nil),
	}

	return evaluator
}

type Result interface {
	fmt.Stringer
}

func (evaluator *Evaluator) Evaluate(statementsChannel <-chan parser.ParsedStatement, bufferSize int) <-chan Result {
	resultsChannel := make(chan Result, bufferSize)

	go evaluator.evaluatingThread(statementsChannel, resultsChannel)

	return resultsChannel
}

func (evaluator *Evaluator) evaluatingThread(statementsChannel <-chan parser.ParsedStatement, resultsChannel chan<- Result) {
	for statement := range statementsChannel {
		if statement.Error != nil {
			resultsChannel <- Error{Err: statement.Error}
			continue
		}

		result, err := evaluator.evaluateStatement(evaluator.globalEnv, statement.Statement)
		if err != nil {
			resultsChannel <- Error{Err: err}
			continue
		}

		resultsChannel <- result
	}

	close(resultsChannel)
}
