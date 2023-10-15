package interpreter

import (
	"banek/ast/expressions"
	"banek/interpreter/objects"
	"banek/tokens"
)

type prefixOperation func(operand objects.Object) (objects.Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperation{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixBangOperation,
}

func (interpreter *Interpreter) evalPrefixOperation(env *environment, expression expressions.PrefixOperation) (objects.Object, error) {
	operand, err := interpreter.evalExpression(env, expression.Operand)
	if err != nil {
		return nil, err
	}

	operation := prefixOperations[expression.Operator.Type]
	if operation == nil {
		return nil, UnknownOperatorError{expression.Operator.Type}
	}

	return operation(operand)
}

func evalPrefixMinusOperation(operand objects.Object) (objects.Object, error) {
	integer, ok := operand.(objects.Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), operand}
	}

	return -integer, nil
}

func evalPrefixBangOperation(operand objects.Object) (objects.Object, error) {
	boolean, ok := operand.(objects.Boolean)
	if !ok {
		return nil, InvalidOperandError{tokens.Bang.String(), operand}
	}

	return !boolean, nil
}
