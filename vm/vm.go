package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
)

const stackSize = 16 * 1024

type vm struct {
	program bytecode.Executable

	opStack [stackSize]objects.Object
	opSP    int

	globalScope  scope
	currentScope *scope
}

type opHandler func(*vm) error

var ops = []opHandler{
	instruction.PushDuplicate:         (*vm).opPushDuplicate,
	instruction.PushConst:             (*vm).opPushConst,
	instruction.PushLocal:             (*vm).opPushLocal,
	instruction.PushGlobal:            (*vm).opPushGlobal,
	instruction.PushCaptured:          (*vm).opPushCaptured,
	instruction.PushBuiltin:           (*vm).opPushBuiltin,
	instruction.PushCollectionElement: (*vm).opPushCollectionElement,

	instruction.Pop:                  (*vm).opPop,
	instruction.PopLocal:             (*vm).opPopLocal,
	instruction.PopGlobal:            (*vm).opPopGlobal,
	instruction.PopCaptured:          (*vm).opPopCaptured,
	instruction.PopCollectionElement: (*vm).opPopCollectionElement,

	instruction.OperationInfix:  (*vm).opInfixOperation,
	instruction.OperationPrefix: (*vm).opPrefixOperation,

	instruction.Branch:        (*vm).opBranch,
	instruction.BranchIfFalse: (*vm).opBranchIfFalse,

	instruction.Call:   (*vm).opCall,
	instruction.Return: (*vm).opReturn,

	instruction.NewArray:    (*vm).opNewArray,
	instruction.NewFunction: (*vm).opNewFunction,
}

func Execute(program bytecode.Executable) error {
	vm := &vm{
		program: program,
		globalScope: scope{
			variables: make([]objects.Object, program.NumGlobals),
			code:      program.Code,
		},
	}
	vm.currentScope = &vm.globalScope

	return vm.run()
}

type ErrUnknownOperation struct {
	Operation instruction.Operation
}

func (err ErrUnknownOperation) Error() string {
	return "unknown operation: " + err.Operation.String()
}

func (vm *vm) run() error {
	for vm.currentScope != nil {
		for vm.hasCode() {
			operation := vm.readOperation()

			if operation == instruction.Invalid || operation >= instruction.Operation(len(ops)) {
				return ErrUnknownOperation{Operation: operation}
			}

			err := ops[operation](vm)
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
