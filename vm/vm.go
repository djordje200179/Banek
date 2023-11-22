package vm

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/vm/scopes"
	"banek/vm/stack"
)

type handler func(program *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack)

var handlers = [...]handler{
	instrs.OpPushDup:       opPushDup,
	instrs.OpPushConst:     opPushConst,
	instrs.OpPush0:         opPush0,
	instrs.OpPush1:         opPush1,
	instrs.OpPush2:         opPush2,
	instrs.OpPushUndefined: opPushUndefined,
	instrs.OpPushBuiltin:   opPushBuiltin,
	instrs.OpPushLocal:     opPushLocal,
	instrs.OpPushLocal0:    opPushLocal0,
	instrs.OpPushLocal1:    opPushLocal1,
	instrs.OpPushGlobal:    opPushGlobal,
	instrs.OpPushCaptured:  opPushCaptured,
	instrs.OpPushCollElem:  opPushCollElem,

	instrs.OpPop:         opPop,
	instrs.OpPopLocal:    opPopLocal,
	instrs.OpPopLocal0:   opPopLocal0,
	instrs.OpPopLocal1:   opPopLocal1,
	instrs.OpPopGlobal:   opPopGlobal,
	instrs.OpPopCaptured: opPopCaptured,
	instrs.OpPopCollElem: opPopCollElem,

	instrs.OpBinaryOp: opBinaryOp,
	instrs.OpUnaryOp:  opUnaryOp,

	instrs.OpBranch:        opBranch,
	instrs.OpBranchIfFalse: opBranchIfFalse,

	instrs.OpCall:   opCall,
	instrs.OpReturn: opReturn,

	instrs.OpNewArray: opNewArray,
	instrs.OpNewFunc:  opNewFunc,
}

func Execute(program bytecode.Executable) {
	scopeStack := scopes.NewStack(program)

	var operandStack stack.Stack

	for {
		opcode := scopeStack.ReadOpcode()

		if opcode == instrs.OpHalt {
			break
		}

		handlers[opcode](&program, scopeStack, &operandStack)
	}
}
