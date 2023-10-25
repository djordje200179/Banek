package vm

import (
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
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

func (vm *vm) opBranch() {
	offset := binary.LittleEndian.Uint16(vm.program.Code[vm.pc+1:])
	vm.pc += int(offset)
}

func (vm *vm) opBranchIfFalse() error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	boolOperand, ok := operand.(objects.Boolean)
	if !ok {
		return errors.ErrInvalidOperand{Operation: instruction.BranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		vm.opBranch()
	}

	return nil
}

func (vm *vm) opNewArray() error {
	size := binary.LittleEndian.Uint16(vm.program.Code[vm.pc+1:])

	array, err := vm.popMany(int(size))
	if err != nil {
		return err
	}

	_ = vm.push(objects.Array(array))

	return nil
}

func (vm *vm) opCollectionAccess() error {
	indexObject, err := vm.pop()
	if err != nil {
		return err
	}

	collectionObject, err := vm.pop()
	if err != nil {
		return err
	}

	var result objects.Object
	switch collection := collectionObject.(type) {
	case objects.Array:
		index, ok := indexObject.(objects.Integer)
		if !ok {
			return errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: indexObject}
		}

		if index < 0 {
			index = objects.Integer(len(collection)) + index
		}

		if index < 0 || index >= objects.Integer(len(collection)) {
			return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(collection)}
		}

		result = collection[index]
	default:
		return errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: indexObject}
	}

	_ = vm.push(result)

	return nil
}
