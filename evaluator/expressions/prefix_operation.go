package expressions

import (
	"banek/evaluator/objects"
	"banek/tokens"
)

type prefixOperationEvaluator func(objects.Object) (objects.Object, error)

var prefixOperations = map[tokens.TokenType]prefixOperationEvaluator{
	tokens.Minus: evalPrefixMinusOperation,
	tokens.Bang:  evalPrefixBangOperation,
}

func evalPrefixOperation(operator tokens.Token, operand objects.Object) (objects.Object, error) {
	operation := prefixOperations[operator.Type]
	if operation == nil {
		return nil, UnknownOperatorError{operator.Type}
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
