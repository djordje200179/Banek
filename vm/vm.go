package vm

import (
	"banek/bytecode"
	"banek/exec/objects"
	"banek/tokens"
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

var infixOperations = map[bytecode.Operation]tokens.TokenType{
	bytecode.Negate:   tokens.Minus,
	bytecode.Add:      tokens.Plus,
	bytecode.Subtract: tokens.Minus,
	bytecode.Multiply: tokens.Asterisk,
	bytecode.Divide:   tokens.Slash,

	bytecode.Equals:           tokens.Equals,
	bytecode.NotEquals:        tokens.NotEquals,
	bytecode.LessThan:         tokens.LessThan,
	bytecode.LessThanOrEquals: tokens.LessThanOrEquals,
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
		case bytecode.Pop:
			err := vm.opPop()
			if err != nil {
				return err
			}
		case bytecode.Add, bytecode.Subtract, bytecode.Multiply, bytecode.Divide,
			bytecode.Equals, bytecode.NotEquals,
			bytecode.LessThan, bytecode.LessThanOrEquals:
			err := vm.opInfixOperation(infixOperations[operation])
			if err != nil {
				return err
			}
		case bytecode.Negate:
			err := vm.opPrefixOperation(infixOperations[operation])
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
