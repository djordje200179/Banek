package opstack

import (
	"banek/runtime/objs"
	"errors"
)

const stackSize = 4 * 1024

type Stack struct {
	arr [stackSize]objs.Obj

	ptr int
}

var ErrStackOverflow = errors.New("stack overflow")

func (s *Stack) Push(obj objs.Obj) {
	if s.ptr >= stackSize {
		panic(ErrStackOverflow)
	}

	s.arr[s.ptr] = obj
	s.ptr++
}

func (s *Stack) Pop() objs.Obj {
	s.ptr--

	elem := s.arr[s.ptr]
	s.arr[s.ptr] = objs.Obj{}

	return elem
}

func (s *Stack) Pop2() (objs.Obj, objs.Obj) {
	s.ptr -= 2

	elem1 := s.arr[s.ptr]
	elem2 := s.arr[s.ptr+1]

	s.arr[s.ptr] = objs.Obj{}
	s.arr[s.ptr+1] = objs.Obj{}

	return elem1, elem2
}

func (s *Stack) Pop3() (objs.Obj, objs.Obj, objs.Obj) {
	s.ptr -= 3

	elem1 := s.arr[s.ptr]
	elem2 := s.arr[s.ptr+1]
	elem3 := s.arr[s.ptr+2]

	s.arr[s.ptr] = objs.Obj{}
	s.arr[s.ptr+1] = objs.Obj{}
	s.arr[s.ptr+2] = objs.Obj{}

	return elem1, elem2, elem3
}

func (s *Stack) PopMany(arr []objs.Obj) {
	nextPtr := s.ptr - len(arr)
	copy(arr, s.arr[nextPtr:s.ptr])
	clear(s.arr[nextPtr:s.ptr])
	s.ptr = nextPtr
}

func (s *Stack) PeekMany(cnt int) []objs.Obj {
	return s.arr[s.ptr-cnt : s.ptr]
}

func (s *Stack) Shrink(cnt int) {
	clear(s.arr[s.ptr-cnt : s.ptr])
	s.ptr -= cnt
}

func (s *Stack) Grow(cnt int) {
	if s.ptr+cnt >= stackSize {
		panic(ErrStackOverflow)
	}

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
	s.arr[s.ptr-1], s.arr[s.ptr-2] = s.arr[s.ptr-2], s.arr[s.ptr-1]
}

func (s *Stack) ReadVar(bp, index int) objs.Obj {
	return s.arr[bp+index]
}

func (s *Stack) WriteVar(bp, index int, obj objs.Obj) {
	s.arr[bp+index] = obj
}

func (s *Stack) SP() int {
	return s.ptr
}
