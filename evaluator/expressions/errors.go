package expressions

import (
	"banek/ast"
	"banek/evaluator/objects"
	"banek/tokens"
	"fmt"
)

type UnknownOperatorError struct {
	Operator tokens.TokenType
}

func (err UnknownOperatorError) Error() string {
	return "unknown operator: " + err.Operator.String()
}

type UnknownExpressionError struct {
	Expression ast.Expression
}

func (err UnknownExpressionError) Error() string {
	return "unknown expression: " + err.Expression.String()
}

type InvalidOperandError struct {
	Operator string
	Operand  objects.Object
}

func (err InvalidOperandError) Error() string {
	return fmt.Sprintf("invalid operand of type %s for operator %s", err.Operand.Type(), err.Operator)
}
