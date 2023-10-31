package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/tokens"
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
		reverseOperands := false
		switch expression.Operator.Type {
		case tokens.GreaterThan, tokens.GreaterThanOrEquals:
			reverseOperands = true
		}

		return compiler.compileInfixOperation(expression, reverseOperands)
	case expressions.PrefixOperation:
		return compiler.compilePrefixOperation(expression)
	case expressions.If:
		// TODO: Implement
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
		err = compiler.compileStatement(expression.Body)
		if err != nil {
			return err
		}
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

		if variableContainerIndex == len(compiler.containerStack)-1 && variableContainerIndex != 0 {
			container.emitInstruction(instruction.PushLocal, variableIndex)
			return nil
		} else if variableContainerIndex == 0 {
			container.emitInstruction(instruction.PushGlobal, variableIndex)
			return nil
		}

		capturedVariableLevel := len(compiler.containerStack) - 1 - variableContainerIndex

		capturedVariableIndex := container.(*functionGenerator).addCapturedVariable(capturedVariableLevel, variableIndex)
		container.emitInstruction(instruction.PushCaptured, capturedVariableIndex)

		return nil
	default:
		return errors.ErrUnknownExpression{Expression: expression}
	}

	return nil
}

func (compiler *compiler) compileInfixOperation(expression expressions.InfixOperation, reverseOperands bool) error {
	var firstOperand, secondOperand ast.Expression
	if reverseOperands {
		firstOperand = expression.Right
		secondOperand = expression.Left
	} else {
		firstOperand = expression.Left
		secondOperand = expression.Right
	}

	err := compiler.compileExpression(firstOperand)
	if err != nil {
		return err
	}

	err = compiler.compileExpression(secondOperand)
	if err != nil {
		return err
	}

	operator := expression.Operator.Type

	container := compiler.topContainer()

	switch operator {
	case tokens.Plus:
		container.emitInstruction(instruction.Add)
	case tokens.Minus:
		container.emitInstruction(instruction.Subtract)
	case tokens.Asterisk:
		container.emitInstruction(instruction.Multiply)
	case tokens.Slash:
		container.emitInstruction(instruction.Divide)
	case tokens.Equals:
		container.emitInstruction(instruction.Equals)
	case tokens.NotEquals:
		container.emitInstruction(instruction.NotEquals)
	case tokens.LessThan, tokens.GreaterThan:
		container.emitInstruction(instruction.LessThan)
	case tokens.LessThanOrEquals, tokens.GreaterThanOrEquals:
		container.emitInstruction(instruction.LessThanOrEquals)
	default:
		return errors.ErrUnknownOperator{Operator: operator}
	}

	return nil
}

func (compiler *compiler) compilePrefixOperation(expression expressions.PrefixOperation) error {
	err := compiler.compileExpression(expression.Operand)
	if err != nil {
		return err
	}

	operator := expression.Operator.Type

	container := compiler.topContainer()

	switch operator {
	case tokens.Minus:
		container.emitInstruction(instruction.Negate)
	case tokens.Bang:
		container.emitInstruction(instruction.Negate)
	default:
		return errors.ErrUnknownOperator{Operator: operator}
	}

	return nil
}
