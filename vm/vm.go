package vm

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/objs"
)

type vm struct {
	program bytecode.Executable

	operandStack
	scopeStack

	halted bool
}

type handler func(vm *vm)

var handlers = [...]handler{
	instrs.OpPushDup:       (*vm).opPushDup,
	instrs.OpPushConst:     (*vm).opPushConst,
	instrs.OpPush0:         (*vm).opPush0,
	instrs.OpPush1:         (*vm).opPush1,
	instrs.OpPush2:         (*vm).opPush2,
	instrs.OpPushUndefined: (*vm).opPushUndefined,
	instrs.OpPushBuiltin:   (*vm).opPushBuiltin,
	instrs.OpPushLocal:     (*vm).opPushLocal,
	instrs.OpPushLocal0:    (*vm).opPushLocal0,
	instrs.OpPushLocal1:    (*vm).opPushLocal1,
	instrs.OpPushGlobal:    (*vm).opPushGlobal,
	instrs.OpPushCaptured:  (*vm).opPushCaptured,
	instrs.OpPushCollElem:  (*vm).opPushCollElem,

	instrs.OpPop:         (*vm).opPop,
	instrs.OpPopLocal:    (*vm).opPopLocal,
	instrs.OpPopLocal0:   (*vm).opPopLocal0,
	instrs.OpPopLocal1:   (*vm).opPopLocal1,
	instrs.OpPopGlobal:   (*vm).opPopGlobal,
	instrs.OpPopCaptured: (*vm).opPopCaptured,
	instrs.OpPopCollElem: (*vm).opPopCollElem,

	instrs.OpBinaryOp: (*vm).opBinaryOp,
	instrs.OpUnaryOp:  (*vm).opUnaryOp,

	instrs.OpBranch:        (*vm).opBranch,
	instrs.OpBranchIfFalse: (*vm).opBranchIfFalse,

	instrs.OpCall:   (*vm).opCall,
	instrs.OpReturn: (*vm).opReturn,
	instrs.OpHalt:   (*vm).opHalt,

	instrs.OpNewArray: (*vm).opNewArray,
	instrs.OpNewFunc:  (*vm).opNewFunc,
}

func Execute(program bytecode.Executable) {
	vm := vm{
		program: program,
		scopeStack: scopeStack{
			globalScope: scope{
				vars: make([]objs.Obj, program.NumGlobals),
				code: program.Code,
			},
		},
	}
	vm.activeScope = vm.globalScope

	for !vm.halted {
		opcode := instrs.Opcode(vm.activeScope.code[vm.activeScope.pc])
		vm.activeScope.pc++

		handlers[opcode](&vm)
	}
}

func (vm *vm) readOperand(width int) int {
	operand := instrs.ReadOperandValue(vm.activeScope.code[vm.activeScope.pc:], width)
	vm.activeScope.pc += width

	return operand
}
