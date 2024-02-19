package stack

import (
	"banek/runtime"
	"errors"
)

const stackSize = 4 * 1024

type Stack struct {
	arr [stackSize]runtime.Obj

	ptr int
}

var ErrStackOverflow = errors.New("stack overflow")

func (stack *Stack) Push(obj runtime.Obj) {
	if stack.ptr >= stackSize {
		panic(ErrStackOverflow)
	}

	stack.arr[stack.ptr] = obj
	stack.ptr++
}

func (stack *Stack) Pop() runtime.Obj {
	stack.ptr--

	elem := stack.arr[stack.ptr]
	stack.arr[stack.ptr] = nil

	return elem
}

func (stack *Stack) PopMany(arr []runtime.Obj) {
	nextPtr := stack.ptr - len(arr)
	copy(arr, stack.arr[nextPtr:stack.ptr])
	clear(stack.arr[nextPtr:stack.ptr])
	stack.ptr = nextPtr
}

func (stack *Stack) Dup() {
	top := stack.arr[stack.ptr-1]
	stack.Push(top)
}

func (stack *Stack) Dup2() {
	top := stack.arr[stack.ptr-1]
	prev := stack.arr[stack.ptr-2]

	stack.Push(prev)
	stack.Push(top)
}

func (stack *Stack) Dup3() {
	top := stack.arr[stack.ptr-1]
	prev := stack.arr[stack.ptr-2]
	prevPrev := stack.arr[stack.ptr-3]

	stack.Push(prevPrev)
	stack.Push(prev)
	stack.Push(top)
}

func (stack *Stack) Swap() {
	top := stack.arr[stack.ptr-1]
	stack.arr[stack.ptr-1] = stack.arr[stack.ptr-2]
	stack.arr[stack.ptr-2] = top
}
