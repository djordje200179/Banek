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

	code   bytecode.Code
	pc     int
	locals []objs.Obj

	halted bool
}

type handler func(vm *vm) error

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

func Execute(program bytecode.Executable) error {
	vm := vm{
		program: program,
		scopeStack: scopeStack{
			globalScope: scope{
				vars: make([]objs.Obj, program.NumGlobals),
				code: program.Code,
			},
		},
		code: program.Code,
	}
	vm.activeScope = &vm.globalScope
	vm.locals = vm.globalScope.vars

	for !vm.halted {
		opcode := instrs.Opcode(vm.code[vm.pc])
		vm.pc++

		err := handlers[opcode](&vm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *vm) readOperand(width int) int {
	rawOperand := vm.code[vm.pc : vm.pc+width]
	operand := instrs.ReadOperandValue(rawOperand, width)

	vm.pc += width

	return operand
}

func (vm *vm) getLocal(index int) objs.Obj {
	return vm.locals[index]
}

func (vm *vm) setLocal(index int, value objs.Obj) {
	vm.locals[index] = value
}
