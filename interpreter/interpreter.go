package interpreter

import (
	"banek/ast"
	"banek/interpreter/environments"
	"banek/interpreter/results"
	"runtime"
)

func Interpret(stmtsChan <-chan ast.Statement, bufferSize int) <-chan results.Result {
	resChan := make(chan results.Result, bufferSize)

	go evalThread(stmtsChan, resChan)

	return resChan
}

type interpreter struct {
	globalEnv *environments.Env
}

func evalThread(stmtsChan <-chan ast.Statement, resChan chan<- results.Result) {
	runtime.LockOSThread()

	interpreter := &interpreter{
		globalEnv: environments.New(nil, 0),
	}

	for stmt := range stmtsChan {
		res, err := interpreter.evalStmt(interpreter.globalEnv, stmt)
		if err != nil {
			resChan <- results.Error{Err: err}
			continue
		}

		resChan <- res
	}

	close(resChan)
}
