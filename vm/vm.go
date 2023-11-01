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
			switch operation {
			case instruction.PushConst:
				err := vm.opPushConst()
				if err != nil {
					return err
				}
			case instruction.PopLocal:
				err := vm.opPopLocal()
				if err != nil {
					return err
				}
			case instruction.PopGlobal:
				err := vm.opPopGlobal()
				if err != nil {
					return err
				}
			case instruction.PopCaptured:
				err := vm.opPopCaptured()
				if err != nil {
					return err
				}
			case instruction.Pop:
				err := vm.opPop()
				if err != nil {
					return err
				}
			case instruction.PushLocal:
				err := vm.opPushLocal()
				if err != nil {
					return err
				}
			case instruction.PushGlobal:
				err := vm.opPushGlobal()
				if err != nil {
					return err
				}
			case instruction.PushBuiltin:
				err := vm.opPushBuiltin()
				if err != nil {
					return err
				}
			case instruction.PushCaptured:
				err := vm.opPushCaptured()
				if err != nil {
					return err
				}
			case instruction.OperationInfix:
				err := vm.opInfixOperation()
				if err != nil {
					return err
				}
			case instruction.OperationPrefix:
				err := vm.opPrefixOperation()
				if err != nil {
					return err
				}
			case instruction.Branch:
				vm.opBranch()
			case instruction.BranchIfFalse:
				err := vm.opBranchIfFalse()
				if err != nil {
					return err
				}
			case instruction.NewArray:
				err := vm.opNewArray()
				if err != nil {
					return err
				}
			case instruction.NewFunction:
				err := vm.opNewFunction()
				if err != nil {
					return err
				}
			case instruction.Call:
				err := vm.opCall()
				if err != nil {
					return err
				}
			case instruction.Return:
				err := vm.opReturn()
				if err != nil {
					return err
				}
			default:
				return ErrUnknownOperation{Operation: operation}
			}
		}

		err := vm.opReturn()
		if err != nil {
			return err
		}
	}

	return nil
}
