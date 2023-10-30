package compiler

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"slices"
)

type functionGenerator struct {
	parameters []string
	locals     []string
	captures   []bytecode.CaptureInfo

	code []byte
}

func (generator *functionGenerator) addParameters(parameters []string) error {
	for i, firstParameter := range generator.parameters {
		for j := i + 1; j < len(generator.parameters); j++ {
			if firstParameter == generator.parameters[j] {
				return errors.ErrIdentifierAlreadyDefined{Identifier: firstParameter}
			}
		}
	}

	generator.parameters = slices.Clone(parameters)
	generator.locals = slices.Clone(parameters)

	return nil
}

func (generator *functionGenerator) addVariable(name string) (int, error) {
	if index := slices.Index(generator.locals, name); index != -1 {
		return 0, errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	generator.locals = append(generator.locals, name)

	return len(generator.locals) - 1, nil
}

func (generator *functionGenerator) addCapturedVariable(index, level int) int {
	captureInfo := bytecode.CaptureInfo{Index: index, Level: level}

	if captureIndex := slices.Index(generator.captures, captureInfo); captureIndex != -1 {
		return captureIndex
	}

	generator.captures = append(generator.captures, captureInfo)

	return len(generator.captures) - 1
}

func (generator *functionGenerator) getVariable(name string) int {
	return slices.Index(generator.locals, name)
}

func (generator *functionGenerator) emitInstruction(operation instruction.Operation, operands ...int) {
	instruction := instruction.MakeInstruction(operation, operands...)

	newCode := make([]byte, len(generator.code)+len(instruction))
	copy(newCode, generator.code)
	copy(newCode[len(generator.code):], instruction)

	generator.code = newCode
}

func (generator *functionGenerator) patchInstructionOperand(address int, operandIndex int, newValue int) {
	operation := instruction.Operation(generator.code[address])
	opInfo := operation.Info()

	instructionCode := generator.code[address : address+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	copy(instructionCode[operandOffset:], instruction.MakeOperandValue(newValue, operandWidth))
}

func (generator *functionGenerator) currentAddress() int {
	return len(generator.code)
}

func (generator *functionGenerator) makeFunction() bytecode.FunctionTemplate {
	return bytecode.FunctionTemplate{
		Code: generator.code,

		Parameters:   generator.parameters,
		NumLocals:    len(generator.locals),
		CapturesInfo: generator.captures,
	}
}

func (generator *functionGenerator) isGlobal() bool {
	return false
}
