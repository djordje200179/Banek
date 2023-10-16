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
