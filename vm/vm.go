package vm

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/types"
)

type vm struct {
	program bytecode.Executable

	operandStack
	scopeStack
}

type handler func(vm *vm, scope *scope) error

var handlers = [...]handler{
	instrs.OpPushDup:      (*vm).opPushDup,
	instrs.OpPushConst:    (*vm).opPushConst,
	instrs.OpPushLocal:    (*vm).opPushLocal,
	instrs.OpPushLocal0:   (*vm).opPushLocal0,
	instrs.OpPushLocal1:   (*vm).opPushLocal1,
	instrs.OpPushGlobal:   (*vm).opPushGlobal,
	instrs.OpPushCaptured: (*vm).opPushCaptured,
	instrs.OpPushCollElem: (*vm).opPushCollElem,

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

	instrs.OpCallFunc:    (*vm).opCallFunc,
	instrs.OpCallBuiltin: (*vm).opCallBuiltin,
	instrs.OpReturn:      (*vm).opReturn,

	instrs.OpNewArray: (*vm).opNewArray,
	instrs.OpNewFunc:  (*vm).opNewFunc,
}

func Execute(program bytecode.Executable) error {
	vm := &vm{
		program: program,
		scopeStack: scopeStack{
			globalScope: scope{
				vars: make([]types.Obj, program.NumGlobals),
				code: program.Code,
			},
		},
	}
	vm.currScope = &vm.globalScope

	return vm.run()
}

func (vm *vm) run() error {
	for {
		scope := vm.currScope

		opcode := scope.readOpcode()
		err := handlers[opcode](vm, scope)
		if err != nil {
			return err
		}
	}
}
