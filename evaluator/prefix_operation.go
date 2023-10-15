package evaluator

import (
	"banek/tokens"
)

type prefixOperation func(operand Object) (Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperation{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixBangOperation,
}

func (evaluator *Evaluator) evalPrefixOperation(operator tokens.Token, operand Object) (Object, error) {
	operation := prefixOperations[operator.Type]
	if operation == nil {
		return nil, UnknownOperatorError{operator.Type}
	}

	return operation(operand)
}

func evalPrefixMinusOperation(operand Object) (Object, error) {
	integer, ok := operand.(Integer)
	if !ok {
		return nil, InvalidOperandError{tokens.Minus.String(), operand}
	}

	return -integer, nil
}

func evalPrefixBangOperation(operand Object) (Object, error) {
	boolean, ok := operand.(Boolean)
	if !ok {
		return nil, InvalidOperandError{tokens.Bang.String(), operand}
	}

	return !boolean, nil
}
