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

	instrs.OpLoadGlobal:   (*emulator).handleLoadGlobal,
	instrs.OpLoadCaptured: nil,
	instrs.OpLoadLocal:    (*emulator).handleLoadLocal,
	instrs.OpLoadLocal0:   (*emulator).handleLoadLocal0,
	instrs.OpLoadLocal1:   (*emulator).handleLoadLocal1,
	instrs.OpLoadLocal2:   (*emulator).handleLoadLocal2,

	instrs.OpConst0:     (*emulator).handleConst0,
	instrs.OpConst1:     (*emulator).handleConst1,
	instrs.OpConst2:     (*emulator).handleConst2,
	instrs.OpConst3:     (*emulator).handleConst3,
	instrs.OpConstN1:    (*emulator).handleConstN1,
	instrs.OpConstInt:   (*emulator).handleConstInt,
	instrs.OpConstStr:   (*emulator).handleConstStr,
	instrs.OpConstTrue:  (*emulator).handleConstTrue,
	instrs.OpConstFalse: (*emulator).handleConstFalse,
	instrs.OpConstUndef: (*emulator).handleConstUndef,
	instrs.OpBuiltin:    (*emulator).handleBuiltin,

	instrs.OpMakeArray: (*emulator).handleMakeArray,
	instrs.OpNewArray:  (*emulator).handleNewArray,
	instrs.OpMakeFunc:  (*emulator).handleMakeFunc,

	instrs.OpPop:           (*emulator).handlePop,
	instrs.OpStoreGlobal:   (*emulator).handleStoreGlobal,
	instrs.OpStoreCaptured: nil,
	instrs.OpStoreLocal:    (*emulator).handleStoreLocal,
	instrs.OpStoreLocal0:   (*emulator).handleStoreLocal0,
	instrs.OpStoreLocal1:   (*emulator).handleStoreLocal1,
	instrs.OpStoreLocal2:   (*emulator).handleStoreLocal2,

	instrs.OpCollGet: (*emulator).handleCollGet,
	instrs.OpCollSet: (*emulator).handleCollSet,

	instrs.OpDup:  (*emulator).handleDup,
	instrs.OpDup2: (*emulator).handleDup2,
	instrs.OpSwap: (*emulator).handleSwap,

	instrs.OpCompEq: (*emulator).handleCompEq,
	instrs.OpCompNe: (*emulator).handleCompNeq,
	instrs.OpCompLt: (*emulator).handleCompLt,
	instrs.OpCompLe: (*emulator).handleCompLe,
	instrs.OpCompGt: (*emulator).handleCompGt,
	instrs.OpCompGe: (*emulator).handleCompGe,

	instrs.OpNeg: (*emulator).handleNeg,
	instrs.OpNot: (*emulator).handleNot,

	instrs.OpAdd: (*emulator).handleAdd,
	instrs.OpSub: (*emulator).handleSub,
	instrs.OpMul: (*emulator).handleMul,
	instrs.OpDiv: (*emulator).handleDiv,
	instrs.OpMod: (*emulator).handleMod,
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
