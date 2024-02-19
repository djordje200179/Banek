package runtime

import "fmt"

type InvalidOperandsError struct {
	Operator BinaryOperator

	Left, Right Obj
}

func (err InvalidOperandsError) Error() string {
	return "invalid operands to " + err.Operator.String() + ": " + err.Left.String() + " and " + err.Right.String()
}

type InvalidOperandError struct {
	Operator UnaryOperator

	Operand Obj
}

func (err InvalidOperandError) Error() string {
	return "invalid operand to " + err.Operator.String() + ": " + err.Operand.String()
}

type TooManyArgsError struct {
	Expected, Actual int
}

func (err TooManyArgsError) Error() string {
	return fmt.Sprintf("too many arguments: expected %d, got %d", err.Expected, err.Actual)
}

type NotCallableError struct {
	Func Obj
}

func (err NotCallableError) Error() string {
	return "not callable: " + err.Func.String()
}

type InvalidTypeError struct {
	BuiltinName string
	ArgIndex    int

	Arg Obj
}

func (err InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type for argument %d to %s: %s", err.ArgIndex, err.BuiltinName, err.Arg)
}
