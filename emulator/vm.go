package emulator

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/emulator/callstack"
	"banek/emulator/opstack"
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
	instrs.OpPushCollElem: (*emulator).handlePushCollElem,
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
	instrs.OpPopCollElem: (*emulator).handlePopCollElem,
	instrs.OpPopLocal:    (*emulator).handlePopLocal,
	instrs.OpPopLocal0:   (*emulator).handlePopLocal0,
	instrs.OpPopLocal1:   (*emulator).handlePopLocal1,
	instrs.OpPopLocal2:   (*emulator).handlePopLocal2,

	instrs.OpDup:  (*emulator).handleDup,
	instrs.OpDup2: (*emulator).handleDup2,
	instrs.OpSwap: (*emulator).handleSwap,

	instrs.OpBinaryAdd: (*emulator).handleBinaryAdd,
	instrs.OpBinarySub: (*emulator).handleBinarySub,
	instrs.OpBinaryMul: (*emulator).handleBinaryMul,
	instrs.OpBinaryDiv: (*emulator).handleBinaryDiv,
	instrs.OpBinaryMod: (*emulator).handleBinaryMod,

	instrs.OpBinaryEq: (*emulator).handleBinaryEq,
	instrs.OpBinaryNe: (*emulator).handleBinaryNeq,
	instrs.OpBinaryLt: (*emulator).handleCompLt,
	instrs.OpBinaryLe: (*emulator).handleCompLe,
	instrs.OpBinaryGt: (*emulator).handleCompGt,
	instrs.OpBinaryGe: (*emulator).handleCompGe,

	instrs.OpUnaryNeg: (*emulator).handleUnaryNeg,
	instrs.OpUnaryNot: (*emulator).handleUnaryNot,

	instrs.OpMakeArray: (*emulator).handleMakeArray,
	instrs.OpNewArray:  (*emulator).handleNewArray,
	instrs.OpMakeFunc:  (*emulator).handleMakeFunc,
}

type emulator struct {
	program bytecode.Executable

	frame callstack.Frame

	opStack   opstack.Stack
	callStack callstack.Stack
}

func (e *emulator) readOpcode() instrs.Opcode {
	opcode := instrs.Opcode(e.program.Code[e.frame.PC])
	e.frame.PC++

	return opcode
}

func (e *emulator) readOperand(op instrs.Opcode, index int) int {
	width := op.Info().Operands[index].Width

	operandSlice := e.program.Code[e.frame.PC : e.frame.PC+width]
	operandValue := instrs.ReadOperandValue(operandSlice)
	e.frame.PC += width

	return operandValue
}

func (e *emulator) movePC(offset int) { e.frame.PC += offset }

func Execute(program bytecode.Executable) (err error) {
	defer func() {
		status := recover()
		if statusErr, ok := status.(error); ok {
			err = statusErr
		}
	}()

	e := emulator{
		program: program,
	}
	e.opStack.Grow(program.FuncPool[0].NumLocals)

	for {
		opcode := e.readOpcode()
		handlers[opcode](&e)
	}
}
