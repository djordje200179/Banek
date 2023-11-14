package errors

import (
	"banek/runtime/types"
	"strconv"
	"strings"
)

type ErrTooManyArgs struct {
	Expected int
	Received int
}

func (err ErrTooManyArgs) Error() string {
	var sb strings.Builder

	sb.WriteString("too many arguments: expected ")
	sb.WriteString(strconv.Itoa(err.Expected))
	sb.WriteString(", received ")
	sb.WriteString(strconv.Itoa(err.Received))

	return sb.String()
}

type ErrNotCallable struct {
	Obj types.Obj
}

func (err ErrNotCallable) Error() string {
	return "not callable: " + err.Obj.String()
}
