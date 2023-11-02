package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/compiler/scopes"
)

func (compiler *compiler) compileStatement(statement ast.Statement) error {
	scope := compiler.topScope()

	switch statement := statement.(type) {
	case statements.Expression:
		err := compiler.compileExpression(statement.Expression)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.Pop)

		return nil
	case statements.If:
		return compiler.compileIfStatement(statement)
	case statements.Block:
		return compiler.compileBlockStatement(statement)
	case statements.Function:
		return compiler.compileFunctionStatement(statement)
	case statements.Return:
		err := compiler.compileExpression(statement.Value)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.Return)

		return nil
	case statements.VariableDeclaration:
		return compiler.compileVariableDeclaration(statement)
	case statements.While:
		return compiler.compileWhileStatement(statement)
	default:
		return ast.ErrUnknownStatement{Statement: statement}
	}
}

func (compiler *compiler) compileIfStatement(statement statements.If) error {
	err := compiler.compileExpression(statement.Condition)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddress := scope.CurrAddr()
	scope.EmitInstr(instruction.BranchIfFalse, 0)

	err = compiler.compileStatement(statement.Consequence)
	if err != nil {
		return err
	}

	elseAddress := scope.CurrAddr()

	if statement.Alternative != nil {
		secondPatchAddress := elseAddress
		scope.EmitInstr(instruction.Branch, 0)
		elseAddress += instruction.Branch.Info().Size()

		err = compiler.compileStatement(statement.Alternative)
		if err != nil {
			return err
		}

		outAddress := scope.CurrAddr()
		scope.PatchInstrOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-instruction.Branch.Info().Size())
	}

	scope.PatchInstrOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

	return nil
}

func (compiler *compiler) compileBlockStatement(statement statements.Block) error {
	scope := compiler.topScope()

	blockScope := &scopes.Block{
		Index:  scope.NextBlockIndex(),
		Parent: scope,
	}

	compiler.pushScope(blockScope)

	for _, statement := range statement.Statements {
		err := compiler.compileStatement(statement)
		if err != nil {
			return err
		}
	}

	compiler.popScope()

	return nil
}

func (compiler *compiler) compileFunctionStatement(statement statements.Function) error {
	functionScope := new(scopes.Function)

	parameterNames := make([]string, len(statement.Parameters))
	for i, parameter := range statement.Parameters {
		parameterNames[i] = parameter.String()
	}

	err := functionScope.AddParams(parameterNames)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	variableIndex, err := scope.AddVar(statement.Name.String(), false)
	if err != nil {
		return err
	}

	compiler.pushScope(functionScope)
	err = compiler.compileStatement(statement.Body)
	if err != nil {
		return err
	}
	compiler.popScope()

	functionTemplate := functionScope.MakeFunction()
	functionTemplate.Name = statement.Name.String()

	functionIndex := compiler.addFunction(functionTemplate)

	if functionTemplate.IsClosure() {
		scope.EmitInstr(instruction.NewFunction, functionIndex)
	} else {
		functionObject := &bytecode.Function{
			TemplateIndex: functionIndex,
		}

		scope.EmitInstr(instruction.PushConst, compiler.addConstant(functionObject))
	}

	if scope.IsGlobal() {
		scope.EmitInstr(instruction.PopGlobal, variableIndex)
	} else {
		scope.EmitInstr(instruction.PopLocal, variableIndex)
	}

	return nil
}

func (compiler *compiler) compileVariableDeclaration(statement statements.VariableDeclaration) error {
	err := compiler.compileExpression(statement.Value)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	index, err := scope.AddVar(statement.Name.String(), statement.Mutable)
	if err != nil {
		return err
	}

	if scope.IsGlobal() {
		scope.EmitInstr(instruction.PopGlobal, index)
	} else {
		scope.EmitInstr(instruction.PopLocal, index)
	}

	return nil
}

func (compiler *compiler) compileWhileStatement(statement statements.While) error {
	scope := compiler.topScope()

	conditionAddress := scope.CurrAddr()

	err := compiler.compileExpression(statement.Condition)
	if err != nil {
		return err
	}

	conditionalBranchAddress := scope.CurrAddr()

	scope.EmitInstr(instruction.BranchIfFalse, 0)

	err = compiler.compileStatement(statement.Body)
	if err != nil {
		return err
	}

	bodyOutAddress := scope.CurrAddr()
	loopOutAddress := bodyOutAddress + instruction.Branch.Info().Size()

	scope.EmitInstr(instruction.Branch, conditionAddress-loopOutAddress)

	scope.PatchInstrOperand(conditionalBranchAddress, 0, loopOutAddress-conditionalBranchAddress-instruction.BranchIfFalse.Info().Size())

	return nil
}
