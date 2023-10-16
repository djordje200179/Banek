package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/tokens"
)

func (compiler *compiler) compileExpression(expression ast.Expression) error {
	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		integer := objects.Integer(expression)
		compiler.emitInstruction(bytecode.PushConst, compiler.addConstant(integer))
		return nil
	case expressions.BooleanLiteral:
		boolean := objects.Boolean(expression)
		compiler.emitInstruction(bytecode.PushConst, compiler.addConstant(boolean))
		return nil
	case expressions.StringLiteral:
		str := objects.String(expression)
		compiler.emitInstruction(bytecode.PushConst, compiler.addConstant(str))
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
		compiler.emitInstruction(bytecode.Add)
	case tokens.Minus:
		compiler.emitInstruction(bytecode.Subtract)
	case tokens.Asterisk:
		compiler.emitInstruction(bytecode.Multiply)
	case tokens.Slash:
		compiler.emitInstruction(bytecode.Divide)
	case tokens.Equals:
		compiler.emitInstruction(bytecode.Equals)
	case tokens.NotEquals:
		compiler.emitInstruction(bytecode.NotEquals)
	case tokens.LessThan, tokens.GreaterThan:
		compiler.emitInstruction(bytecode.LessThan)
	case tokens.LessThanOrEquals, tokens.GreaterThanOrEquals:
		compiler.emitInstruction(bytecode.LessThanOrEquals)
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
		compiler.emitInstruction(bytecode.Negate)
	case tokens.Bang:
		compiler.emitInstruction(bytecode.Negate)
	default:
		return errors.UnknownOperatorError{Operator: operator}
	}

	return nil
}
