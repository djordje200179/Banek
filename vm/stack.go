package vm

import (
	"banek/exec/objects"
)

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
}

func (vm *vm) push(object objects.Object) error {
	if vm.sp >= stackSize {
		return ErrStackOverflow{}
	}

	vm.stack[vm.sp] = object
	vm.sp++

	return nil
}

func (vm *vm) pop() (objects.Object, error) {
	if vm.sp <= 0 {
		return nil, ErrStackOverflow{}
	}

	vm.sp--
	object := vm.stack[vm.sp]

	return object, nil
}

func (vm *vm) popMany(count int) ([]objects.Object, error) {
	if vm.sp < count {
		return nil, ErrStackOverflow{}
	}

	elements := make([]objects.Object, count)
	for i := 0; i < count; i++ {
		object, err := vm.pop()
		if err != nil {
			return nil, err
		}

		elements[i] = object
	}

	return elements, nil
}
