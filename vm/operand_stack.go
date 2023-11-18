package vm

import (
	"banek/runtime/objs"
)

const stackSize = 4 * 1024

type operandStack struct {
	array [stackSize]objs.Obj

	ptr int
}

type ErrStackOverflow struct{}

func (err ErrStackOverflow) Error() string {
	return "stack overflow"
}

func (stack *operandStack) peek() objs.Obj {
	return stack.array[stack.ptr-1]
}

func (stack *operandStack) swap(obj objs.Obj) {
	stack.array[stack.ptr-1] = obj
}

func (stack *operandStack) push(obj objs.Obj) error {
	if stack.ptr >= stackSize {
		return ErrStackOverflow{}
	}

	stack.array[stack.ptr] = obj
	stack.ptr++

	return nil
}

func (stack *operandStack) pop() objs.Obj {
	stack.ptr--

	elem := stack.array[stack.ptr]
	stack.array[stack.ptr] = objs.Obj{}

	return elem
}

func (stack *operandStack) popMany(arr []objs.Obj) {
	nextPtr := stack.ptr - len(arr)
	copy(arr, stack.array[nextPtr:stack.ptr])

	for i := nextPtr; i < stack.ptr; i++ {
		stack.array[i] = objs.Obj{}
	}

	stack.ptr = nextPtr
}
