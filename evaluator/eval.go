package evaluator

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/evaluator/expressions"
	"banek/evaluator/objects"
	"banek/parser"
	"fmt"
)

func evalStatement(statement ast.Statement) (objects.Object, error) {
	switch statement := statement.(type) {
	case statements.Expression:
		return expressions.EvalExpression(statement.Expression)
	case statements.If:
		condition, err := expressions.EvalExpression(statement.Condition)
		if err != nil {
			return nil, err
		}

		if condition == objects.Boolean(true) {
			return evalStatement(statement.Consequence)
		} else if statement.Alternative != nil {
			return evalStatement(statement.Alternative)
		} else {
			return objects.None{}, nil
		}
	case statements.Block:
		for _, statement := range statement.Statements {
			_, err := evalStatement(statement)
			if err != nil {
				return nil, err
			}
		}

		return objects.None{}, nil
	default:
		return nil, nil
	}
}

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

		result, err := evalStatement(statement.Statement)
		if err != nil {
			fmt.Println(err)
			continue
		}

		objectChannel <- result
	}

	close(objectChannel)
}
