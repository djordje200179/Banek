package vm

import (
	"banek/exec/operations"
	"banek/tokens"
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

func (vm *vm) opPop() error {
	_, err := vm.pop()
	return err
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

	result, err := operations.EvalInfixOperation(left, right, tokens.Plus)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}

func (vm *vm) opSubtract() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalInfixOperation(left, right, tokens.Minus)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}

func (vm *vm) opMultiply() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalInfixOperation(left, right, tokens.Asterisk)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}

func (vm *vm) opDivide() error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalInfixOperation(left, right, tokens.Slash)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}
