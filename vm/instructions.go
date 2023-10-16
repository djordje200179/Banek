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

func (vm *vm) opInfixOperation(operation tokens.TokenType) error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalInfixOperation(left, right, operation)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}

func (vm *vm) opPrefixOperation(operation tokens.TokenType) error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalPrefixOperation(operand, operation)
	if err != nil {
		return err
	}

	_ = vm.push(result)

	return nil
}
