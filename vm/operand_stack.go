package vm

import "banek/runtime/types"

const stackSize = 4 * 1024

type operandStack struct {
	array [stackSize]types.Obj

	ptr int
}

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
}

func (stack *operandStack) peek() types.Obj {
	return stack.array[stack.ptr-1]
}

func (stack *operandStack) push(object types.Obj) error {
	if stack.ptr >= stackSize {
		return ErrStackOverflow{}
	}

	stack.array[stack.ptr] = object
	stack.ptr++

	return nil
}

func (stack *operandStack) popOne() types.Obj {
	stack.ptr--

	elem := stack.array[stack.ptr]
	stack.array[stack.ptr] = nil

	return elem
}

func (stack *operandStack) popMany(arr []types.Obj) {
	nextPtr := stack.ptr - len(arr)
	copy(arr, stack.array[nextPtr:stack.ptr])

	for i := nextPtr; i < stack.ptr; i++ {
		stack.array[i] = nil
	}

	stack.ptr = nextPtr
}
