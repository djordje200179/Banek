package unaryops

import (
	"banek/runtime/objs"
	"fmt"
)

var unaryOps = [unaryOperatorCount][objs.TypeCount]func(objs.Obj) (objs.Obj, bool){
	NegOperator: {
		objs.Int: negateInt,
	},
	NotOperator: {
		objs.Bool: invertBool,
	},
}

type InvalidOperandError struct {
	Operator UnaryOperator

	Operand objs.Obj
}

func (err InvalidOperandError) Error() string {
	return fmt.Sprintf("invalid operand for %s: %s", err.Operator.String(), err.Operand.String())
}

func (op UnaryOperator) Eval(o objs.Obj) (objs.Obj, error) {
	err := InvalidOperandError{op, o}

	handler := unaryOps[op][o.Type]
	if handler == nil {
		return objs.Obj{}, err
	}

	result, ok := handler(o)
	if !ok {
		return objs.Obj{}, err
	}

	return result, nil
}
