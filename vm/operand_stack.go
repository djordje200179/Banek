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

type ErrStackUnderflow struct{}

func (err ErrStackUnderflow) Error() string {
	return "stack underflow"
}

func (stack *operandStack) peek() (types.Obj, error) {
	if stack.ptr <= 0 {
		return nil, ErrStackUnderflow{}
	}

	return stack.array[stack.ptr-1], nil
}

func (stack *operandStack) push(object types.Obj) error {
	if stack.ptr >= stackSize {
		return ErrStackOverflow{}
	}

	stack.array[stack.ptr] = object
	stack.ptr++

	return nil
}

func (stack *operandStack) popOne() (types.Obj, error) {
	if stack.ptr <= 0 {
		return nil, ErrStackUnderflow{}
	}

	stack.ptr--

	elem := stack.array[stack.ptr]
	stack.array[stack.ptr] = nil

	return elem, nil
}

func (stack *operandStack) popTwo() (types.Obj, types.Obj, error) {
	if stack.ptr <= 1 {
		return nil, nil, ErrStackUnderflow{}
	}

	stack.ptr -= 2

	firstElem := stack.array[stack.ptr]
	secondElem := stack.array[stack.ptr+1]

	stack.array[stack.ptr] = nil
	stack.array[stack.ptr+1] = nil

	return firstElem, secondElem, nil
}

func (stack *operandStack) popMany(arr []types.Obj) error {
	if stack.ptr < len(arr) {
		return ErrStackUnderflow{}
	}

	nextPtr := stack.ptr - len(arr)
	copy(arr, stack.array[nextPtr:stack.ptr])

	for i := nextPtr; i < stack.ptr; i++ {
		stack.array[i] = nil
	}

	stack.ptr = nextPtr

	return nil
}

func (stack *operandStack) decreaseSP(size int) error {
	if stack.ptr < size {
		return ErrStackUnderflow{}
	}

	stack.ptr -= size

	for i := 0; i < size; i++ {
		stack.array[stack.ptr+i] = nil
	}

	return nil
}
