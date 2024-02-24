package binaryops

import (
	"banek/runtime/objs"
	"fmt"
)

var binaryOps = [binaryOperatorCount][objs.TypeCount][objs.TypeCount]func(objs.Obj, objs.Obj) objs.Obj{
	AddOperator: {
		objs.Int: {
			objs.Int: addInts,
		},
		objs.String: {
			objs.String: addStrings,
		},
		objs.Array: {
			objs.Array: concatArrays,
		},
	},
	SubOperator: {
		objs.Int: {
			objs.Int: subInts,
		},
	},
	MulOperator: {
		objs.Int: {
			objs.Int: mulInts,
		},
		objs.String: {
			objs.Int: repeatStrings,
		},
		objs.Array: {
			objs.Int: repeatArray,
		},
	},
	DivOperator: {
		objs.Int: {
			objs.Int: divInts,
		},
	},
	ModOperator: {
		objs.Int: {
			objs.Int: modInts,
		},
	},
}

type InvalidOperandsError struct {
	Operator BinaryOperator

	Left, Right objs.Obj
}

func (err InvalidOperandsError) Error() string {
	return fmt.Sprintf("invalid operands for %s: %s and %s", err.Operator.String(), err.Left.String(), err.Right.String())
}

func (op BinaryOperator) Eval(left, right objs.Obj) objs.Obj {
	handler := binaryOps[op][left.Type()][right.Type()]
	if handler == nil {
		panic(InvalidOperandsError{op, left, right})
	}

	return handler(left, right)
}
