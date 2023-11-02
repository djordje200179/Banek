package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/compiler/scopes"
	"banek/exec/errors"
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
		err := compiler.compileExpression(statement.Condition)
		if err != nil {
			return err
		}

		firstPatchAddress := scope.CurrAddr()
		scope.EmitInstr(instruction.BranchIfFalse, 0)

		err = compiler.compileStatement(statement.Consequence)
		if err != nil {
			return err
		}

		elseAddress := scope.CurrAddr()

		if statement.Alternative != nil {
			branchSize := instruction.Branch.Info().Size()

			secondPatchAddress := elseAddress
			scope.EmitInstr(instruction.Branch, 0)
			elseAddress += branchSize

			err = compiler.compileStatement(statement.Alternative)
			if err != nil {
				return err
			}

			outAddress := scope.CurrAddr()
			scope.PatchInstrOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-branchSize)
		}

		scope.PatchInstrOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

		return nil
	case statements.Block:
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
	case statements.Function:
		functionScope := new(scopes.Function)

		parameterNames := make([]string, len(statement.Parameters))
		for i, parameter := range statement.Parameters {
			parameterNames[i] = parameter.String()
		}

		err := functionScope.AddParams(parameterNames)
		if err != nil {
			return err
		}

		variableIndex, err := scope.AddVar(statement.Name.String(), false)
		if err != nil {
			return err
		}

		compiler.scopes = append(compiler.scopes, functionScope)
		err = compiler.compileStatement(statement.Body)
		if err != nil {
			return err
		}
		compiler.scopes = compiler.scopes[:len(compiler.scopes)-1]

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
	case statements.Return:
		err := compiler.compileExpression(statement.Value)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.Return)

		return nil
	case statements.VariableDeclaration:
		err := compiler.compileExpression(statement.Value)
		if err != nil {
			return err
		}

		index, err := scope.AddVar(statement.Name.String(), !statement.Const)
		if err != nil {
			return err
		}

		if scope.IsGlobal() {
			scope.EmitInstr(instruction.PopGlobal, index)
		} else {
			scope.EmitInstr(instruction.PopLocal, index)
		}

		return nil
	case statements.While:
		// TODO: Implement
	default:
		return errors.ErrUnknownStatement{Statement: statement}
	}

	return nil
}
