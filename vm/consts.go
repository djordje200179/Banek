package vm

import (
	"banek/exec/objects"
	"strconv"
	"strings"
)

type ErrConstNotDefined struct {
	Index int
}

func (err ErrConstNotDefined) Error() string {
	var sb strings.Builder

	sb.WriteString("undefined constant at index ")
	sb.WriteString(strconv.Itoa(err.Index))

	return sb.String()
}

func (vm *vm) getConst(index int) (objects.Object, error) {
	if index >= len(vm.program.ConstsPool) {
		return nil, ErrConstNotDefined{index}
	}

	return vm.program.ConstsPool[index], nil
}
