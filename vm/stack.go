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
	if count == 0 {
		return nil, nil
	} else if count == 1 {
		object, err := vm.pop()
		if err != nil {
			return nil, err
		}

		return []objects.Object{object}, nil
	}

	if vm.opSP < count {
		return nil, ErrStackOverflow{}
	}

	nextSP := vm.opSP - count

	elements := vm.opStack[nextSP:vm.opSP]
	vm.opSP = nextSP

	return slices.Clone(elements), nil
}
