package compiler

import (
	"banek/ast"
	"banek/ast/stmts"
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/compiler/scopes"
)

func (compiler *compiler) compileStmt(stmt ast.Stmt) error {
	scope := compiler.topScope()

	switch stmt := stmt.(type) {
	case stmts.Expr:
		err := compiler.compileExpr(stmt.Expr)
		if err != nil {
			return err
		}

		scope.EmitInstr(instrs.OpPop)

		return nil
	case stmts.If:
		return compiler.compileIfStatement(stmt)
	case stmts.Block:
		return compiler.compileBlock(stmt)
	case stmts.Func:
		return compiler.compileFuncStmt(stmt)
	case stmts.Return:
		err := compiler.compileExpr(stmt.Value)
		if err != nil {
			return err
		}

		scope.EmitInstr(instrs.OpReturn)

		return nil
	case stmts.VarDecl:
		return compiler.compileVarDeclaration(stmt)
	case stmts.While:
		return compiler.compileWhile(stmt)
	default:
		return ast.ErrUnknownStmt{Stmt: stmt}
	}
}

func (compiler *compiler) compileIfStatement(stmt stmts.If) error {
	err := compiler.compileExpr(stmt.Cond)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instrs.OpBranchIfFalse, 0)

	err = compiler.compileStmt(stmt.Consequence)
	if err != nil {
		return err
	}

	elseAddr := scope.CurrAddr()

	branchSize := instrs.OpBranch.Info().Size()
	branchIfFalseSize := instrs.OpBranchIfFalse.Info().Size()

	if stmt.Alternative != nil {
		secondPatchAddr := elseAddr
		scope.EmitInstr(instrs.OpBranch, 0)
		elseAddr += branchSize

		err = compiler.compileStmt(stmt.Alternative)
		if err != nil {
			return err
		}

		outAddr := scope.CurrAddr()
		scope.PatchInstrOperand(secondPatchAddr, 0, outAddr-secondPatchAddr-branchSize)
	}

	scope.PatchInstrOperand(firstPatchAddr, 0, elseAddr-firstPatchAddr-branchIfFalseSize)

	return nil
}

func (compiler *compiler) compileBlock(stmt stmts.Block) error {
	scope := compiler.topScope()

	blockScope := &scopes.Block{
		Index:  scope.NextBlockIndex(),
		Parent: scope,
	}

	compiler.pushScope(blockScope)

	for _, stmt := range stmt.Stmts {
		err := compiler.compileStmt(stmt)
		if err != nil {
			return err
		}
	}

	compiler.popScope()

	return nil
}

func (compiler *compiler) compileFuncStmt(stmt stmts.Func) error {
	funcScope := new(scopes.Function)

	paramNames := make([]string, len(stmt.Params))
	for i, param := range stmt.Params {
		paramNames[i] = param.String()
	}

	err := funcScope.AddParams(paramNames)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	varIndex, err := scope.AddVar(stmt.Name.String(), false)
	if err != nil {
		return err
	}

	compiler.pushScope(funcScope)
	err = compiler.compileStmt(stmt.Body)
	if err != nil {
		return err
	}
	compiler.popScope()

	funcTemplate := funcScope.MakeFunction()
	funcTemplate.Name = stmt.Name.String()

	funcIndex := compiler.addFunc(funcTemplate)

	if funcTemplate.IsClosure() {
		scope.EmitInstr(instrs.OpNewFunc, funcIndex)
	} else {
		functionObject := &bytecode.Func{
			TemplateIndex: funcIndex,
		}

		scope.EmitInstr(instrs.OpPushConst, compiler.addConst(functionObject))
	}

	if scope.IsGlobal() {
		scope.EmitInstr(instrs.OpPopGlobal, varIndex)
	} else {
		scope.EmitInstr(instrs.OpPopLocal, varIndex)
	}

	return nil
}

func (compiler *compiler) compileVarDeclaration(stmt stmts.VarDecl) error {
	err := compiler.compileExpr(stmt.Value)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	index, err := scope.AddVar(stmt.Name.String(), stmt.Mutable)
	if err != nil {
		return err
	}

	if scope.IsGlobal() {
		scope.EmitInstr(instrs.OpPopGlobal, index)
	} else {
		scope.EmitInstr(instrs.OpPopLocal, index)
	}

	return nil
}

func (compiler *compiler) compileWhile(stmt stmts.While) error {
	scope := compiler.topScope()

	condAddr := scope.CurrAddr()

	err := compiler.compileExpr(stmt.Cond)
	if err != nil {
		return err
	}

	condBranchAddr := scope.CurrAddr()

	scope.EmitInstr(instrs.OpBranchIfFalse, 0)

	err = compiler.compileStmt(stmt.Body)
	if err != nil {
		return err
	}

	bodyOutAddr := scope.CurrAddr()
	loopOutAddr := bodyOutAddr + instrs.OpBranch.Info().Size()

	scope.EmitInstr(instrs.OpBranch, condAddr-loopOutAddr)

	scope.PatchInstrOperand(condBranchAddr, 0, loopOutAddr-condBranchAddr-instrs.OpBranchIfFalse.Info().Size())

	return nil
}
