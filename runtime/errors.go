package runtime

import "fmt"

type InvalidOperandsError struct {
	Operator BinaryOperator

	Left, Right Obj
}

func (err InvalidOperandsError) Error() string {
	return fmt.Sprintf("invalid operands for %s: %s and %s", err.Operator.String(), err.Left.String(), err.Right.String())
}

type InvalidOperandError struct {
	Operator UnaryOperator

	Operand Obj
}

func (err InvalidOperandError) Error() string {
	return fmt.Sprintf("invalid operand for %s: %s", err.Operator.String(), err.Operand.String())
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
	return fmt.Sprintf("invalid type for %d. argument of %s: %s", err.ArgIndex+1, err.BuiltinName, err.Arg.String())
}

type NotIndexableError struct {
	Coll Obj
	Key  Obj
}

func (err NotIndexableError) Error() string {
	return fmt.Sprintf("not indexable: %s, for key: %s", err.Coll.String(), err.Key.String())
}
