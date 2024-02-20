package runtime

import (
	"banek/runtime/objs"
	"fmt"
)

type TooManyArgsError struct {
	Expected, Actual int
}

func (err TooManyArgsError) Error() string {
	return fmt.Sprintf("too many arguments: expected %d, got %d", err.Expected, err.Actual)
}

type NotCallableError struct {
	Func objs.Obj
}

func (err NotCallableError) Error() string {
	return "not callable: " + err.Func.String()
}
