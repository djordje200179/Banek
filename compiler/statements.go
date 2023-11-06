package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/compiler/scopes"
)

func (compiler *compiler) compileStmt(stmt ast.Statement) error {
	scope := compiler.topScope()

	switch stmt := stmt.(type) {
	case statements.Expr:
		err := compiler.compileExpr(stmt.Expr)
		if err != nil {
			return err
		}

		scope.EmitInstr(instructions.OpPop)

		return nil
	case statements.If:
		return compiler.compileIfStatement(stmt)
	case statements.Block:
		return compiler.compileBlock(stmt)
	case statements.Func:
		return compiler.compileFuncStmt(stmt)
	case statements.Return:
		err := compiler.compileExpr(stmt.Value)
		if err != nil {
			return err
		}

		scope.EmitInstr(instructions.OpReturn)

		return nil
	case statements.VarDecl:
		return compiler.compileVarDeclaration(stmt)
	case statements.While:
		return compiler.compileWhile(stmt)
	default:
		return ast.ErrUnknownStmt{Stmt: stmt}
	}
}

func (compiler *compiler) compileIfStatement(stmt statements.If) error {
	err := compiler.compileExpr(stmt.Cond)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instructions.OpBranchIfFalse, 0)

	err = compiler.compileStmt(stmt.Consequence)
	if err != nil {
		return err
	}

	elseAddr := scope.CurrAddr()

	branchSize := instructions.OpBranch.Info().Size()
	branchIfFalseSize := instructions.OpBranchIfFalse.Info().Size()

	if stmt.Alternative != nil {
		secondPatchAddr := elseAddr
		scope.EmitInstr(instructions.OpBranch, 0)
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

func (compiler *compiler) compileBlock(stmt statements.Block) error {
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

func (compiler *compiler) compileFuncStmt(stmt statements.Func) error {
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
		scope.EmitInstr(instructions.OpNewFunc, funcIndex)
	} else {
		functionObject := &bytecode.Func{
			TemplateIndex: funcIndex,
		}

		scope.EmitInstr(instructions.OpPushConst, compiler.addConst(functionObject))
	}

	if scope.IsGlobal() {
		scope.EmitInstr(instructions.OpPopGlobal, varIndex)
	} else {
		scope.EmitInstr(instructions.OpPopLocal, varIndex)
	}

	return nil
}

func (compiler *compiler) compileVarDeclaration(stmt statements.VarDecl) error {
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
		scope.EmitInstr(instructions.OpPopGlobal, index)
	} else {
		scope.EmitInstr(instructions.OpPopLocal, index)
	}

	return nil
}

func (compiler *compiler) compileWhile(stmt statements.While) error {
	scope := compiler.topScope()

	condAddr := scope.CurrAddr()

	err := compiler.compileExpr(stmt.Cond)
	if err != nil {
		return err
	}

	condBranchAddr := scope.CurrAddr()

	scope.EmitInstr(instructions.OpBranchIfFalse, 0)

	err = compiler.compileStmt(stmt.Body)
	if err != nil {
		return err
	}

	bodyOutAddr := scope.CurrAddr()
	loopOutAddr := bodyOutAddr + instructions.OpBranch.Info().Size()

	scope.EmitInstr(instructions.OpBranch, condAddr-loopOutAddr)

	scope.PatchInstrOperand(condBranchAddr, 0, loopOutAddr-condBranchAddr-instructions.OpBranchIfFalse.Info().Size())

	return nil
}
