package errors

import "banek/runtime/objs"

type ErrUnknownOperator struct {
	Operator string
}

func (err ErrUnknownOperator) Error() string {
	return "unknown operator: " + err.Operator
}

type ErrNotBool struct {
	Obj objs.Obj
}

func (err ErrNotBool) Error() string {
	return "not bool: " + err.Obj.String()
}
