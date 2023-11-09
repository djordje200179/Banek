package builtins

import (
	"strconv"
	"strings"
)

type ErrIncorrectArgNum struct {
	Expected int
	Got      int
}

func (err ErrIncorrectArgNum) Error() string {
	var sb strings.Builder

	sb.WriteString("incorrect number of arguments: expected ")
	sb.WriteString(strconv.Itoa(err.Expected))
	sb.WriteString(", got ")
	sb.WriteString(strconv.Itoa(err.Got))

	return sb.String()
}
