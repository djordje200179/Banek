package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/objects"
)

const stackSize = 4 * 1024

type vm struct {
	program bytecode.Executable

	opStack [stackSize]objects.Object
	opSP    int

	globalScope scope
	currScope   *scope
}

type opHandler func(*vm) error

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
		globalScope: scope{
			vars: make([]objects.Object, program.NumGlobals),
			code: program.Code,
		},
	}
	vm.currScope = &vm.globalScope

	return vm.run()
}

type ErrUnknownInstr struct {
	InstrType instructions.Opcode
}

func (err ErrUnknownInstr) Error() string {
	return "unknown instruction " + err.InstrType.String()
}

func (vm *vm) run() error {
	for vm.currScope != nil {
		for vm.hasCode() {
			opcode := vm.readOpcode()
			if opcode == instructions.OpInvalid || opcode >= instructions.Opcode(len(ops)) {
				return ErrUnknownInstr{InstrType: opcode}
			}

			err := ops[opcode](vm)
			if err != nil {
				return err
			}
		}

		err := vm.opReturn()
		if err != nil {
			return err
		}
	}

	return nil
}
