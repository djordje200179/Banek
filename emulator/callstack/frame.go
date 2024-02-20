package callstack

import (
	"banek/emulator/function"
)

type Frame struct {
	PC int
	BP int

	Func *function.Obj
}
