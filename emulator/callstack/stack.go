package callstack

import (
	"banek/emulator/function"
	"errors"
)

const stackSize = 4 * 1024

type Stack struct {
	arr [stackSize]Frame

	ptr int
}

var ErrStackOverflow = errors.New("recursion limit exceeded")

func (s *Stack) Push(pc, bp int, funcObj *function.Obj) *Frame {
	if s.ptr == stackSize-1 {
		panic(ErrStackOverflow)
	}

	s.ptr++
	s.arr[s.ptr] = Frame{pc, bp, funcObj}

	return &s.arr[s.ptr]
}

func (s *Stack) Pop() {
	s.ptr--
}

func (s *Stack) GlobalFrame() *Frame {
	return &s.arr[0]
}

func (s *Stack) ActiveFrame() *Frame {
	return &s.arr[s.ptr]
}
