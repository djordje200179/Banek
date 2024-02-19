package emulator

import (
	"banek/runtime"
	"errors"
)

const stackSize = 4 * 1024

type stack struct {
	arr [stackSize]runtime.Obj

	ptr int
}

var ErrStackOverflow = errors.New("stack overflow")

func (s *stack) push(obj runtime.Obj) {
	if s.ptr >= stackSize {
		panic(ErrStackOverflow)
	}

	s.arr[s.ptr] = obj
	s.ptr++
}

func (s *stack) pop() runtime.Obj {
	s.ptr--

	elem := s.arr[s.ptr]
	s.arr[s.ptr] = nil

	return elem
}

func (s *stack) popMany(arr []runtime.Obj) {
	nextPtr := s.ptr - len(arr)
	copy(arr, s.arr[nextPtr:s.ptr])
	clear(s.arr[nextPtr:s.ptr])
	s.ptr = nextPtr
}

func (s *stack) dup() {
	top := s.arr[s.ptr-1]
	s.push(top)
}

func (s *stack) dup2() {
	top := s.arr[s.ptr-1]
	prev := s.arr[s.ptr-2]

	s.push(prev)
	s.push(top)
}

func (s *stack) Swap() {
	top := s.arr[s.ptr-1]
	s.arr[s.ptr-1] = s.arr[s.ptr-2]
	s.arr[s.ptr-2] = top
}
