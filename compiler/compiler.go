package compiler

import (
	"banek/ast"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/objects"
	"slices"
)

type codeContainer interface {
	addVariable(name string) (int, error)
	getVariable(name string) int

	emitInstruction(operation instruction.Operation, operands ...int)
	patchInstructionOperand(address int, operandIndex int, newValue int)
	currentAddress() int
}

type compiler struct {
	constants []objects.Object
	functions []bytecode.FunctionTemplate

	containerStack []codeContainer
}

func Compile(statementsChannel <-chan ast.Statement) (bytecode.Executable, error) {
	compiler := &compiler{
		containerStack: []codeContainer{new(executableGenerator)},
	}

	for statement := range statementsChannel {
		err := compiler.compileStatement(statement)
		if err != nil {
			return bytecode.Executable{}, err
		}
	}

	return compiler.makeExecutable(), nil
}

func (compiler *compiler) addConstant(object objects.Object) int {
	if index := slices.Index(compiler.constants, object); index != -1 {
		return index
	}

	index := len(compiler.constants)
	compiler.constants = append(compiler.constants, object)

	return index
}

func (compiler *compiler) addFunction(template bytecode.FunctionTemplate) int {
	index := len(compiler.functions)
	compiler.functions = append(compiler.functions, template)

	return index
}

func (compiler *compiler) topContainer() codeContainer {
	return compiler.containerStack[len(compiler.containerStack)-1]
}

func (compiler *compiler) makeExecutable() bytecode.Executable {
	executableGenerator := compiler.containerStack[0].(*executableGenerator)

	executable := executableGenerator.makeExecutable()
	executable.ConstantsPool = compiler.constants
	executable.FunctionsPool = compiler.functions

	return executable
}
