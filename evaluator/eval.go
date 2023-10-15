package evaluator

import (
	"banek/evaluator/objects"
	"banek/evaluator/statements"
	"banek/parser"
	"fmt"
)

func EvalStatements(statementChannel <-chan parser.ParsedStatement, bufferSize int) <-chan objects.Object {
	objectChannel := make(chan objects.Object, bufferSize)

	go evalThread(statementChannel, objectChannel)

	return objectChannel
}

func evalThread(statementChannel <-chan parser.ParsedStatement, objectChannel chan<- objects.Object) {
	for statement := range statementChannel {
		if statement.Error != nil {
			fmt.Println(statement.Error)
			continue
		}

		result, err := statements.EvalStatement(statement.Statement)
		if err != nil {
			fmt.Println(err)
			continue
		}

		objectChannel <- result
	}

	close(objectChannel)
}
