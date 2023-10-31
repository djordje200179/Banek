package vm

import (
	"banek/exec/objects"
	"fmt"
)

type ErrConstantNotDefined struct {
	Index int
}

func (err ErrConstantNotDefined) Error() string {
	return fmt.Sprintf("undefined constant at Index %d", err.Index)
}

func (vm *vm) getConstant(index int) (objects.Object, error) {
	if index >= len(vm.program.ConstantsPool) {
		return nil, ErrConstantNotDefined{index}
	}

	return vm.program.ConstantsPool[index], nil
}
