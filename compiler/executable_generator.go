package compiler

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"slices"
)

type executableGenerator struct {
	globalVariables []string

	code []byte
}

func (generator *executableGenerator) addVariable(name string) (int, error) {
	if index := slices.Index(generator.globalVariables, name); index != -1 {
		return 0, errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	generator.globalVariables = append(generator.globalVariables, name)

	return len(generator.globalVariables) - 1, nil
}

func (generator *executableGenerator) getVariable(name string) int {
	return slices.Index(generator.globalVariables, name)
}

func (generator *executableGenerator) emitInstruction(operation instruction.Operation, operands ...int) {
	instruction := instruction.MakeInstruction(operation, operands...)

	newCode := make([]byte, len(generator.code)+len(instruction))
	copy(newCode, generator.code)
	copy(newCode[len(generator.code):], instruction)

	generator.code = newCode
}

func (generator *executableGenerator) patchInstructionOperand(address int, operandIndex int, newValue int) {
	operation := instruction.Operation(generator.code[address])
	opInfo := operation.Info()

	instructionCode := generator.code[address : address+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	copy(instructionCode[operandOffset:], instruction.MakeOperandValue(newValue, operandWidth))
}

func (generator *executableGenerator) currentAddress() int {
	return len(generator.code)
}

func (generator *executableGenerator) makeExecutable() bytecode.Executable {
	return bytecode.Executable{
		Code: generator.code,

		NumGlobals: len(generator.globalVariables),
	}
}

func (generator *executableGenerator) isGlobal() bool {
	return true
}
