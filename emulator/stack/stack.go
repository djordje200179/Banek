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

func (s *Stack) Push(obj runtime.Obj) {
	if s.ptr >= stackSize {
		panic(ErrStackOverflow)
	}

	s.arr[s.ptr] = obj
	s.ptr++
}

func (s *Stack) Pop() runtime.Obj {
	s.ptr--

	elem := s.arr[s.ptr]
	s.arr[s.ptr] = nil

	return elem
}

func (s *Stack) PopMany(arr []runtime.Obj) {
	nextPtr := s.ptr - len(arr)
	copy(arr, s.arr[nextPtr:s.ptr])
	clear(s.arr[nextPtr:s.ptr])
	s.ptr = nextPtr
}

func (s *Stack) Grow(cnt int) {
	s.ptr += cnt
}

func (s *Stack) Dup() {
	top := s.arr[s.ptr-1]
	s.Push(top)
}

func (s *Stack) Dup2() {
	top := s.arr[s.ptr-1]
	prev := s.arr[s.ptr-2]

	s.Push(prev)
	s.Push(top)
}

func (s *Stack) Swap() {
	top := s.arr[s.ptr-1]
	s.arr[s.ptr-1] = s.arr[s.ptr-2]
	s.arr[s.ptr-2] = top
}

func (s *Stack) ReadVar(bp, index int) runtime.Obj {
	return s.arr[bp+index]
}

func (s *Stack) WriteVar(bp, index int, obj runtime.Obj) {
	s.arr[bp+index] = obj
}

func (s *Stack) PatchReturn(bp int) {
	s.arr[bp] = s.arr[s.ptr-1]
	s.ptr = bp + 1
}

func (s *Stack) SP() int {
	return s.ptr
}
