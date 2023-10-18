package vm

import (
	"banek/exec/objects"
	"fmt"
)

type ErrConstantNotDefined struct {
	Index uint16
}

func (err ErrConstantNotDefined) Error() string {
	return fmt.Sprintf("undefined constant at index %d", err.Index)
}

func (vm *vm) getConstant(index uint16) (objects.Object, error) {
	if int(index) >= len(vm.program.ConstantsPool) {
		return nil, ErrConstantNotDefined{index}
	}

	return vm.program.ConstantsPool[index], nil
}
