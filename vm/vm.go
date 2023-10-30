package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
	"banek/tokens"
)

const stackSize = 2048

type vm struct {
	program bytecode.Executable

	operationStack        [stackSize]objects.Object
	operationStackPointer int

	globalScope scope
	scopeStack  []*scope
}

func Execute(program bytecode.Executable) error {
	vm := &vm{
		program: program,
		globalScope: scope{
			variables: make([]objects.Object, program.NumGlobals),
			code:      program.Code,
		},
		scopeStack: make([]*scope, 1),
	}
	vm.scopeStack[0] = &vm.globalScope

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
	for len(vm.scopeStack) > 0 {
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

			if operation != instruction.Call {
				vm.movePC(operation.Info().Size())
			}
		}

		err := vm.opReturn()
		if err != nil {
			return err
		}
	}

	return nil
}
