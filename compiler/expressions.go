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
		default:
			return errors.UnknownOperatorError{Operator: operator}
		}

		return nil
	default:
		return errors.UnknownExpressionError{Expression: expression}
	}
}
