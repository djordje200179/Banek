package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
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

type ErrUnknownOperation struct {
	Operation instruction.Operation
}

func (err ErrUnknownOperation) Error() string {
	return "unknown operation: " + err.Operation.String()
}

var infixOperations = map[instruction.Operation]tokens.TokenType{
	instruction.Negate:   tokens.Minus,
	instruction.Add:      tokens.Plus,
	instruction.Subtract: tokens.Minus,
	instruction.Multiply: tokens.Asterisk,
	instruction.Divide:   tokens.Slash,

	instruction.Equals:           tokens.Equals,
	instruction.NotEquals:        tokens.NotEquals,
	instruction.LessThan:         tokens.LessThan,
	instruction.LessThanOrEquals: tokens.LessThanOrEquals,
}

func (vm *vm) run() error {
	for vm.pc = 0; vm.pc < len(vm.program.Code); {
		operation := instruction.Operation(vm.program.Code[vm.pc])
		switch operation {
		case instruction.PushConst:
			err := vm.opPushConst()
			if err != nil {
				return err
			}
		case instruction.Pop:
			err := vm.opPop()
			if err != nil {
				return err
			}
		case instruction.Add, instruction.Subtract, instruction.Multiply, instruction.Divide,
			instruction.Equals, instruction.NotEquals,
			instruction.LessThan, instruction.LessThanOrEquals:
			err := vm.opInfixOperation(infixOperations[operation])
			if err != nil {
				return err
			}
		case instruction.Negate:
			err := vm.opPrefixOperation(infixOperations[operation])
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
		default:
			return ErrUnknownOperation{Operation: operation}
		}

		vm.pc += operation.Info().Size()
	}

	return nil
}
