package analyzer

import (
	"banek/ast"
	"banek/symtable"
	"runtime"
)

func Analyze(stmtChan <-chan ast.Stmt, bufferSize int) <-chan ast.Stmt {
	analyzedStmtChan := make(chan ast.Stmt, bufferSize)

	go analyzingThread(stmtChan, analyzedStmtChan)

	return analyzedStmtChan
}

type analyzer struct {
	symTable *symtable.Table
	loopCnt  int
	funcCnt  int
}

func analyzingThread(stmtChan <-chan ast.Stmt, analyzedStmtChan chan<- ast.Stmt) {
	runtime.LockOSThread()

	analyzer := &analyzer{symTable: symtable.New()}

	for stmt := range stmtChan {
		analyzedStmt, err := analyzer.analyzeStmt(stmt)
		if err != nil {
			close(analyzedStmtChan)
			panic(err)
		}

		analyzedStmtChan <- analyzedStmt
	}

	close(analyzedStmtChan)
}
