package callstack

import (
	"errors"
)

const stackSize = 4 * 1024

type Stack struct {
	arr [stackSize]Frame

	ptr int
}

var ErrStackOverflow = errors.New("recursion limit exceeded")

func (s *Stack) Push(frame Frame) {
	if s.ptr == stackSize-1 {
		panic(ErrStackOverflow)
	}

	s.arr[s.ptr] = frame
	s.ptr++
}

func (s *Stack) Pop() Frame {
	s.ptr--

	return s.arr[s.ptr]
}

//func (s *Stack) GlobalFrame() *Frame {
//	return &s.arr[0]
//}
//
//func (s *Stack) ActiveFrame() *Frame {
//	return &s.arr[s.ptr]
//}
