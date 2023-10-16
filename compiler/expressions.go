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
		err := compiler.compileExpression(expression.Left)
		if err != nil {
			return err
		}

		err = compiler.compileExpression(expression.Right)
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
		default:
			return errors.UnknownOperatorError{Operator: operator}
		}

		return nil
	default:
		return errors.UnknownExpressionError{Expression: expression}
	}
}
