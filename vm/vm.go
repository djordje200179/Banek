package vm

import (
	"banek/bytecode"
	"banek/exec/objects"
)

const stackSize = 2048

type vm struct {
	program *bytecode.Executable

	stack [stackSize]objects.Object
	sp    int

	pc int
}

func Execute(program *bytecode.Executable) error {
	vm := &vm{
		program: program,
	}

	return vm.run()
}

type UnknownOperationError struct {
	Operation bytecode.Operation
}

func (err UnknownOperationError) Error() string {
	return "unknown operation: " + err.Operation.String()
}

func (vm *vm) run() error {
	for vm.pc = 0; vm.pc < len(vm.program.Code); vm.pc++ {
		operation := bytecode.Operation(vm.program.Code[vm.pc])
		switch operation {
		case bytecode.PushConst:
			err := vm.opPushConst()
			if err != nil {
				return err
			}
		case bytecode.Add:
			err := vm.opAdd()
			if err != nil {
				return err
			}
		default:
			return UnknownOperationError{Operation: operation}
		}

		vm.pc += operation.Info().OperandsSize()
	}

	return nil
}
