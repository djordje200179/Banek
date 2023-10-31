package compiler

import (
	"banek/ast"
	"banek/ast/statements"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
)

func (compiler *compiler) compileStatement(statement ast.Statement) error {
	container := compiler.topContainer()

	switch statement := statement.(type) {
	case statements.Expression:
		err := compiler.compileExpression(statement.Expression)
		if err != nil {
			return err
		}

		container.emitInstruction(instruction.Pop)

		return nil
	case statements.If:
		err := compiler.compileExpression(statement.Condition)
		if err != nil {
			return err
		}

		firstPatchAddress := container.currentAddress()
		container.emitInstruction(instruction.BranchIfFalse, 0)

		err = compiler.compileStatement(statement.Consequence)
		if err != nil {
			return err
		}

		elseAddress := container.currentAddress()

		if statement.Alternative != nil {
			branchSize := instruction.Branch.Info().Size()

			secondPatchAddress := elseAddress
			container.emitInstruction(instruction.Branch, 0)
			elseAddress += branchSize

			err = compiler.compileStatement(statement.Alternative)
			if err != nil {
				return err
			}

			outAddress := container.currentAddress()
			container.patchInstructionOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-branchSize)
		}

		container.patchInstructionOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

		return nil
	case statements.Block:
		for _, statement := range statement.Statements {
			err := compiler.compileStatement(statement)
			if err != nil {
				return err
			}
		}

		return nil
	case statements.Function:
		functionGenerator := new(functionGenerator)

		parameterNames := make([]string, len(statement.Parameters))
		for i, parameter := range statement.Parameters {
			parameterNames[i] = parameter.String()
		}

		err := functionGenerator.addParameters(parameterNames)
		if err != nil {
			return err
		}

		variableIndex, err := container.addVariable(statement.Name.String())
		if err != nil {
			return err
		}

		compiler.containerStack = append(compiler.containerStack, functionGenerator)
		err = compiler.compileStatement(statement.Body)
		if err != nil {
			return err
		}
		compiler.containerStack = compiler.containerStack[:len(compiler.containerStack)-1]

		functionTemplate := functionGenerator.makeFunction()
		functionTemplate.Name = statement.Name.String()

		functionIndex := compiler.addFunction(functionTemplate)

		if functionTemplate.IsClosure() {
			container.emitInstruction(instruction.NewFunction, functionIndex)
		} else {
			functionObject := &bytecode.Function{
				TemplateIndex: functionIndex,
			}

			container.emitInstruction(instruction.PushConst, compiler.addConstant(functionObject))
		}

		if len(compiler.containerStack) == 1 {
			container.emitInstruction(instruction.PopGlobal, variableIndex)
		} else {
			container.emitInstruction(instruction.PopLocal, variableIndex)
		}

		return nil
	case statements.Return:
		err := compiler.compileExpression(statement.Value)
		if err != nil {
			return err
		}

		container.emitInstruction(instruction.Return)

		return nil
	case statements.VariableDeclaration:
		err := compiler.compileExpression(statement.Value)
		if err != nil {
			return err
		}

		index, err := container.addVariable(statement.Name.String())
		if err != nil {
			return err
		}

		if len(compiler.containerStack) == 1 {
			container.emitInstruction(instruction.PopGlobal, index)
		} else {
			container.emitInstruction(instruction.PopLocal, index)
		}

		return nil
	case statements.While:
		// TODO: Implement
	default:
		return errors.ErrUnknownStatement{Statement: statement}
	}

	return nil
}
