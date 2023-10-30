package vm

import (
	"banek/exec/objects"
)

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
}

func (vm *vm) push(object objects.Object) error {
	if vm.operationStackPointer >= stackSize {
		return ErrStackOverflow{}
	}

	vm.operationStack[vm.operationStackPointer] = object
	vm.operationStackPointer++

	return nil
}

func (vm *vm) pop() (objects.Object, error) {
	if vm.operationStackPointer <= 0 {
		return nil, ErrStackOverflow{}
	}

	vm.operationStackPointer--
	object := vm.operationStack[vm.operationStackPointer]

	return object, nil
}

func (vm *vm) popMany(count int) ([]objects.Object, error) {
	if vm.operationStackPointer < count {
		return nil, ErrStackOverflow{}
	}

	nextStackPointer := vm.operationStackPointer - count

	elements := vm.operationStack[nextStackPointer:vm.operationStackPointer]
	vm.operationStackPointer = nextStackPointer

	return elements, nil
}
