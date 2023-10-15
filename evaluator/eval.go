package evaluator

import (
	"banek/evaluator/environment"
	"banek/evaluator/objects"
	"banek/evaluator/statements"
	"banek/parser"
	"fmt"
)

func EvalStatements(statementChannel <-chan parser.ParsedStatement, bufferSize int) <-chan objects.Object {
	objectChannel := make(chan objects.Object, bufferSize)
	env := environment.New(nil)

	go evalThread(env, statementChannel, objectChannel)

	return objectChannel
}

func evalThread(env *environment.Environment, statementChannel <-chan parser.ParsedStatement, objectChannel chan<- objects.Object) {
	for statement := range statementChannel {
		if statement.Error != nil {
			fmt.Println(statement.Error)
			continue
		}

		result, err := statements.EvalStatement(env, statement.Statement)
		if err != nil {
			fmt.Println(err)
			continue
		}

		objectChannel <- result
	}

	close(objectChannel)
}
