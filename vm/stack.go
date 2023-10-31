package vm

import (
	"banek/exec/objects"
	"slices"
)

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
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
		return nil, ErrStackOverflow{}
	}

	vm.opSP--
	object := vm.opStack[vm.opSP]

	return object, nil
}

func (vm *vm) popMany(count int) ([]objects.Object, error) {
	if vm.opSP < count {
		return nil, ErrStackOverflow{}
	}

	nextStackPointer := vm.opSP - count

	elements := vm.opStack[nextStackPointer:vm.opSP]
	vm.opSP = nextStackPointer

	return slices.Clone(elements), nil
}
