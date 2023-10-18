package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/tokens"
)

func (compiler *compiler) compileExpression(expression ast.Expression) error {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		integer := objects.Integer(expression)
		compiler.emitInstruction(instruction.PushConst, compiler.addConstant(integer))
		return nil
	case expressions.BooleanLiteral:
		boolean := objects.Boolean(expression)
		compiler.emitInstruction(instruction.PushConst, compiler.addConstant(boolean))
		return nil
	case expressions.StringLiteral:
		str := objects.String(expression)
		compiler.emitInstruction(instruction.PushConst, compiler.addConstant(str))
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
		err := compiler.compileExpression(expression.Condition)
		if err != nil {
			return err
		}

		return nil

		// TODO: Add support for else branch
	case expressions.ArrayLiteral:
		for _, element := range expression {
			err := compiler.compileExpression(element)
			if err != nil {
				return err
			}
		}

		compiler.emitInstruction(instruction.NewArray, len(expression))

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

		compiler.emitInstruction(instruction.CollectionAccess)

		return nil
	default:
		return errors.UnknownExpressionError{Expression: expression}
	}
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

	switch operator {
	case tokens.Plus:
		compiler.emitInstruction(instruction.Add)
	case tokens.Minus:
		compiler.emitInstruction(instruction.Subtract)
	case tokens.Asterisk:
		compiler.emitInstruction(instruction.Multiply)
	case tokens.Slash:
		compiler.emitInstruction(instruction.Divide)
	case tokens.Equals:
		compiler.emitInstruction(instruction.Equals)
	case tokens.NotEquals:
		compiler.emitInstruction(instruction.NotEquals)
	case tokens.LessThan, tokens.GreaterThan:
		compiler.emitInstruction(instruction.LessThan)
	case tokens.LessThanOrEquals, tokens.GreaterThanOrEquals:
		compiler.emitInstruction(instruction.LessThanOrEquals)
	default:
		return errors.UnknownOperatorError{Operator: operator}
	}

	return nil
}

func (compiler *compiler) compilePrefixOperation(expression expressions.PrefixOperation) error {
	err := compiler.compileExpression(expression.Operand)
	if err != nil {
		return err
	}

	operator := expression.Operator.Type

	switch operator {
	case tokens.Minus:
		compiler.emitInstruction(instruction.Negate)
	case tokens.Bang:
		compiler.emitInstruction(instruction.Negate)
	default:
		return errors.UnknownOperatorError{Operator: operator}
	}

	return nil
}
