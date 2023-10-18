package vm

import (
	"banek/exec/objects"
)

type StackOverflowError struct{}

func (err StackOverflowError) Error() string {
	return "stack overflow"
}

func (vm *vm) push(object objects.Object) error {
	if vm.sp >= stackSize {
		return StackOverflowError{}
	}

	vm.stack[vm.sp] = object
	vm.sp++

	return nil
}

func (vm *vm) pop() (objects.Object, error) {
	if vm.sp <= 0 {
		return nil, StackOverflowError{}
	}

	vm.sp--
	object := vm.stack[vm.sp]

	return object, nil
}

func (vm *vm) popMany(count int) ([]objects.Object, error) {
	if vm.sp < count {
		return nil, StackOverflowError{}
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
