package unaryops

import (
	"banek/runtime/objs"
	"fmt"
)

var unaryOps = [unaryOperatorCount][objs.TypeCount]func(objs.Obj) objs.Obj{
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

func (op UnaryOperator) Eval(o objs.Obj) objs.Obj {
	handler := unaryOps[op][o.Type()]
	if handler == nil {
		panic(InvalidOperandError{op, o})
	}

	return handler(o)

}
