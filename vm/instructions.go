package vm

import (
	"banek/exec/objects"
	"encoding/binary"
)

func (vm *vm) opPushConst() error {
	constIndex := binary.LittleEndian.Uint16(vm.program.Code[vm.pc+1:])

	constant, err := vm.getConstant(constIndex)
	if err != nil {
		return err
	}

	err = vm.push(constant)
	if err != nil {
		return err
	}

	return nil
}

func (vm *vm) opAdd() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	leftValue := left.(objects.Integer)
	rightValue := right.(objects.Integer)

	result := leftValue + rightValue
	_ = vm.push(result)

	return nil
}
