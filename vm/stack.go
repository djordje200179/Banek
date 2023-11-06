package vm

import (
	"banek/exec/objects"
)

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
}

type ErrStackUnderflow struct{}

func (err ErrStackUnderflow) Error() string {
	return "stack underflow"
}

func (vm *vm) peek() (objects.Object, error) {
	if vm.opSP <= 0 {
		return nil, ErrStackUnderflow{}
	}

	return vm.opStack[vm.opSP-1], nil
}

func (vm *vm) push(object objects.Object) error {
	if vm.opSP >= stackSize {
		return ErrStackOverflow{}
	}

	vm.opStack[vm.opSP] = object
	vm.opSP++

	return nil
}

func (vm *vm) pop() (objects.Object, error) {
	if vm.opSP <= 0 {
		return nil, ErrStackUnderflow{}
	}

	vm.opSP--
	object := vm.opStack[vm.opSP]

	return object, nil
}

func (vm *vm) popMany(arr []objects.Object) error {
	if vm.opSP < len(arr) {
		return ErrStackUnderflow{}
	}

	nextSP := vm.opSP - len(arr)
	copy(arr, vm.opStack[nextSP:vm.opSP])

	vm.opSP = nextSP

	return nil
}
