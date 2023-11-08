package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
)

type vm struct {
	program bytecode.Executable

	operandStack
	scopeStack
}

type opHandler func(vm *vm, scope *scope) error

var ops = [...]opHandler{
	instructions.OpPushDup:      (*vm).opPushDup,
	instructions.OpPushConst:    (*vm).opPushConst,
	instructions.OpPushLocal:    (*vm).opPushLocal,
	instructions.OpPushGlobal:   (*vm).opPushGlobal,
	instructions.OpPushCaptured: (*vm).opPushCaptured,
	instructions.OpPushBuiltin:  (*vm).opPushBuiltin,
	instructions.OpPushCollElem: (*vm).opPushCollElem,

	instructions.OpPop:         (*vm).opPop,
	instructions.OpPopLocal:    (*vm).opPopLocal,
	instructions.OpPopGlobal:   (*vm).opPopGlobal,
	instructions.OpPopCaptured: (*vm).opPopCaptured,
	instructions.OpPopCollElem: (*vm).opPopCollElem,

	instructions.OpBinaryOp: (*vm).opBinaryOp,
	instructions.OpUnaryOp:  (*vm).opUnaryOp,

	instructions.OpBranch:        (*vm).opBranch,
	instructions.OpBranchIfFalse: (*vm).opBranchIfFalse,

	instructions.OpCall:   (*vm).opCall,
	instructions.OpReturn: (*vm).opReturn,

	instructions.OpNewArray: (*vm).opNewArray,
	instructions.OpNewFunc:  (*vm).opNewFunc,
}

func Execute(program bytecode.Executable) error {
	vm := &vm{
		program: program,
		scopeStack: scopeStack{
			globalScope: scope{
				vars: make([]objects.Object, program.NumGlobals),
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
		err := ops[opcode](vm, scope)
		if err != nil {
			return err
		}

		scope = vm.currScope

		if !scope.hasCode() {
			if scope == &vm.globalScope {
				return nil
			}

			vm.popScope()
		}
	}
}
