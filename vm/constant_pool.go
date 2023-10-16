package vm

import (
	"banek/exec/objects"
	"fmt"
)

type UndefinedConstantError struct {
	Index uint16
}

func (err UndefinedConstantError) Error() string {
	return fmt.Sprintf("undefined constant at index %d", err.Index)
}

func (vm *vm) getConstant(index uint16) (objects.Object, error) {
	if int(index) >= len(vm.program.ConstantsPool) {
		return nil, UndefinedConstantError{index}
	}

	return vm.program.ConstantsPool[index], nil
}
