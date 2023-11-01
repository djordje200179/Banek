package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
)

func (compiler *compiler) compileExpression(expression ast.Expression) error {
	container := compiler.topContainer()

	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		integer := objects.Integer(expression)
		container.emitInstruction(instruction.PushConst, compiler.addConstant(integer))
		return nil
	case expressions.BooleanLiteral:
		boolean := objects.Boolean(expression)
		container.emitInstruction(instruction.PushConst, compiler.addConstant(boolean))
		return nil
	case expressions.StringLiteral:
		str := objects.String(expression)
		container.emitInstruction(instruction.PushConst, compiler.addConstant(str))
		return nil
	case expressions.InfixOperation:
		return compiler.compileInfixOperation(expression)
	case expressions.PrefixOperation:
		return compiler.compilePrefixOperation(expression)
	case expressions.If:
		err := compiler.compileExpression(expression.Condition)
		if err != nil {
			return err
		}

		firstPatchAddress := container.currentAddress()
		container.emitInstruction(instruction.BranchIfFalse, 0)

		err = compiler.compileExpression(expression.Consequence)
		if err != nil {
			return err
		}

		elseAddress := container.currentAddress()

		branchSize := instruction.Branch.Info().Size()

		secondPatchAddress := elseAddress
		container.emitInstruction(instruction.Branch, 0)
		elseAddress += branchSize

		err = compiler.compileExpression(expression.Alternative)
		if err != nil {
			return err
		}

		outAddress := container.currentAddress()
		container.patchInstructionOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-branchSize)

		container.patchInstructionOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

		return nil
	case expressions.ArrayLiteral:
		for _, element := range expression {
			err := compiler.compileExpression(element)
			if err != nil {
				return err
			}
		}

		container.emitInstruction(instruction.NewArray, len(expression))

		return nil
	case expressions.CollectionAccess:
		err := compiler.compileExpression(expression.Collection)
		if err != nil {
			return err
		}

		err = compiler.compileExpression(expression.Key)
		if err != nil {
			return err
		}

		container.emitInstruction(instruction.CollectionAccess)

		return nil
	case expressions.Assignment:
		// TODO: Implement
	case expressions.FunctionCall:
		for _, argument := range expression.Arguments {
			err := compiler.compileExpression(argument)
			if err != nil {
				return err
			}
		}

		err := compiler.compileExpression(expression.Function)
		if err != nil {
			return err
		}

		container.emitInstruction(instruction.Call, len(expression.Arguments))

		return nil
	case expressions.FunctionLiteral:
		functionGenerator := new(functionGenerator)

		parameterNames := make([]string, len(expression.Parameters))
		for i, parameter := range expression.Parameters {
			parameterNames[i] = parameter.String()
		}

		err := functionGenerator.addParameters(parameterNames)
		if err != nil {
			return err
		}

		compiler.containerStack = append(compiler.containerStack, functionGenerator)
		err = compiler.compileExpression(expression.Body)
		if err != nil {
			return err
		}
		functionGenerator.emitInstruction(instruction.Return)
		compiler.containerStack = compiler.containerStack[:len(compiler.containerStack)-1]

		functionTemplate := functionGenerator.makeFunction()

		functionIndex := compiler.addFunction(functionTemplate)

		if functionTemplate.IsClosure() {
			container.emitInstruction(instruction.NewFunction, functionIndex)
		} else {
			functionObject := &bytecode.Function{
				TemplateIndex: functionIndex,
			}

			container.emitInstruction(instruction.PushConst, compiler.addConstant(functionObject))
		}

		return nil
	case expressions.Identifier:
		variableName := expression.String()

		if index := objects.BuiltinFindIndex(variableName); index != -1 {
			container.emitInstruction(instruction.PushBuiltin, index)
			return nil
		}

		var variableContainer codeContainer
		var variableIndex, variableContainerIndex int
		for i := len(compiler.containerStack) - 1; i >= 0; i-- {
			index := compiler.containerStack[i].getVariable(variableName)
			if index == -1 {
				continue
			}

			variableContainer = compiler.containerStack[i]
			variableContainerIndex = i
			variableIndex = index

			break
		}

		if variableContainer == nil {
			return errors.ErrIdentifierNotDefined{Identifier: variableName}
		}

		if variableContainerIndex == 0 {
			container.emitInstruction(instruction.PushGlobal, variableIndex)
			return nil
		} else if variableContainerIndex == len(compiler.containerStack)-1 {
			container.emitInstruction(instruction.PushLocal, variableIndex)
			return nil
		}

		capturedVariableLevel := len(compiler.containerStack) - 2 - variableContainerIndex

		capturedVariableIndex := container.(*functionGenerator).addCapturedVariable(capturedVariableLevel, variableIndex)
		container.emitInstruction(instruction.PushCaptured, capturedVariableIndex)

		return nil
	default:
		return errors.ErrUnknownExpression{Expression: expression}
	}

	return nil
}
