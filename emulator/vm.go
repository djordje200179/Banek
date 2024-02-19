package emulator

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/emulator/scopes"
	"banek/emulator/stack"
	"banek/runtime"
)

var handlers = [...]func(e *emulator){
	instrs.OpHalt:   (*emulator).handleHalt,
	instrs.OpCall:   (*emulator).handleCall,
	instrs.OpReturn: (*emulator).handleReturn,

	instrs.OpJump:        (*emulator).handleJump,
	instrs.OpBranchFalse: (*emulator).handleBranchFalse,
	instrs.OpBranchTrue:  (*emulator).handleBranchTrue,

	instrs.OpPushBuiltin:  (*emulator).handlePushBuiltin,
	instrs.OpPushGlobal:   (*emulator).handlePushGlobal,
	instrs.OpPushCaptured: nil,
	instrs.OpPushCollElem: nil,
	instrs.OpPushLocal:    (*emulator).handlePushLocal,
	instrs.OpPushLocal0:   (*emulator).handlePushLocal0,
	instrs.OpPushLocal1:   (*emulator).handlePushLocal1,
	instrs.OpPushLocal2:   (*emulator).handlePushLocal2,

	instrs.OpPush0:     (*emulator).handlePush0,
	instrs.OpPush1:     (*emulator).handlePush1,
	instrs.OpPush2:     (*emulator).handlePush2,
	instrs.OpPush3:     (*emulator).handlePush3,
	instrs.OpPushN1:    (*emulator).handlePushN1,
	instrs.OpPushInt:   (*emulator).handlePushInt,
	instrs.OpPushStr:   (*emulator).handlePushStr,
	instrs.OpPushTrue:  (*emulator).handlePushTrue,
	instrs.OpPushFalse: (*emulator).handlePushFalse,
	instrs.OpPushUndef: (*emulator).handlePushUndef,

	instrs.OpPop:         (*emulator).handlePop,
	instrs.OpPopGlobal:   (*emulator).handlePopGlobal,
	instrs.OpPopCaptured: nil,
	instrs.OpPopCollElem: nil,
	instrs.OpPopLocal:    (*emulator).handlePopLocal,
	instrs.OpPopLocal0:   (*emulator).handlePopLocal0,
	instrs.OpPopLocal1:   (*emulator).handlePopLocal1,
	instrs.OpPopLocal2:   (*emulator).handlePopLocal2,

	instrs.OpDup:  (*emulator).handleDup,
	instrs.OpDup2: (*emulator).handleDup2,
	instrs.OpDup3: (*emulator).handleDup3,
	instrs.OpSwap: (*emulator).handleSwap,

	instrs.OpBinaryAdd: (*emulator).handleBinaryAdd,
	instrs.OpBinarySub: (*emulator).handleBinarySub,
	instrs.OpBinaryMul: (*emulator).handleBinaryMul,
	instrs.OpBinaryDiv: (*emulator).handleBinaryDiv,
	instrs.OpBinaryMod: (*emulator).handleBinaryMod,

	instrs.OpBinaryEq: (*emulator).handleBinaryEq,
	instrs.OpBinaryNe: (*emulator).handleBinaryNeq,
	instrs.OpBinaryLt: makeComparisonHandler(runtime.LtOperator),
	instrs.OpBinaryLe: makeComparisonHandler(runtime.LtEqOperator),
	instrs.OpBinaryGt: makeComparisonHandler(runtime.GtOperator),
	instrs.OpBinaryGe: makeComparisonHandler(runtime.GtEqOperator),

	instrs.OpUnaryNeg: (*emulator).handleUnaryNeg,
	instrs.OpUnaryNot: (*emulator).handleUnaryNot,

	instrs.OpMakeArray: (*emulator).handleMakeArray,
	instrs.OpNewArray:  (*emulator).handleNewArray,
	instrs.OpMakeFunc:  (*emulator).handleMakeFunc,
}

type emulator struct {
	program *bytecode.Executable

	scopeStack   *scopes.Stack
	operandStack stack.Stack
}

func Execute(program *bytecode.Executable) (err error) {
	defer func() {
		status := recover()
		if statusErr, ok := status.(error); ok {
			err = statusErr
		}
	}()

	e := emulator{
		program:    program,
		scopeStack: scopes.NewStack(program),
	}

	for {
		opcode := e.scopeStack.ReadOpcode()
		handlers[opcode](&e)
	}
}
