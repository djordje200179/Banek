package vm

import (
	"banek/exec/objects"
)

const stackSize = 4 * 1024

type operandStack struct {
	array [stackSize]objects.Object

	ptr int
}

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "operandStack overflow"
}

type ErrStackUnderflow struct{}

func (err ErrStackUnderflow) Error() string {
	return "operandStack underflow"
}

func (stack *operandStack) peek() (objects.Object, error) {
	if stack.ptr <= 0 {
		return nil, ErrStackUnderflow{}
	}

	return stack.array[stack.ptr-1], nil
}

func (stack *operandStack) push(object objects.Object) error {
	if stack.ptr >= stackSize {
		return ErrStackOverflow{}
	}

	stack.array[stack.ptr] = object
	stack.ptr++

	return nil
}

func (stack *operandStack) pop() (objects.Object, error) {
	if stack.ptr <= 0 {
		return nil, ErrStackUnderflow{}
	}

	stack.ptr--

	return stack.array[stack.ptr], nil
}

func (stack *operandStack) popMany(arr []objects.Object) error {
	if stack.ptr < len(arr) {
		return ErrStackUnderflow{}
	}

	nextPtr := stack.ptr - len(arr)
	copy(arr, stack.array[nextPtr:stack.ptr])

	stack.ptr = nextPtr

	return nil
}
