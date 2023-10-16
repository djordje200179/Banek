package compiler

import (
	"banek/ast"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
)

type compiler struct {
	executable *bytecode.Executable
}

func Compile(statementsChannel <-chan ast.Statement) (*bytecode.Executable, error) {
	compiler := &compiler{
		executable: new(bytecode.Executable),
	}

	for statement := range statementsChannel {
		err := compiler.compileStatement(statement)
		if err != nil {
			return nil, err
		}
	}

	return compiler.executable, nil
}

func (compiler *compiler) addConstant(object objects.Object) int {
	for index, constant := range compiler.executable.ConstantsPool {
		if constant == object {
			return index
		}
	}

	index := len(compiler.executable.ConstantsPool)
	compiler.executable.ConstantsPool = append(compiler.executable.ConstantsPool, object)
	return index
}

func (compiler *compiler) emitInstruction(operation instruction.Operation, operands ...int) {
	instr := instruction.MakeInstruction(operation, operands...)
	compiler.executable.Code = append(compiler.executable.Code, instr...)
}

func (compiler *compiler) patchInstructionOperand(address int, operandIndex int, newValue int) {
	operation := instruction.Operation(compiler.executable.Code[address])
	opInfo := operation.Info()

	instructionCode := compiler.executable.Code[address : address+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	copy(instructionCode[operandOffset:], instruction.MakeOperandValue(newValue, operandWidth))
}

func (compiler *compiler) currentAddress() int {
	return len(compiler.executable.Code)
}
