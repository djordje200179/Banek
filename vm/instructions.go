package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
	"banek/tokens"
	"encoding/binary"
	"slices"
)

func (vm *vm) opPushConst() error {
	constIndex := binary.LittleEndian.Uint16(vm.readCode())

	constant, err := vm.getConstant(constIndex)
	if err != nil {
		return err
	}

	return vm.push(constant)
}

func (vm *vm) opPushLocal() error {
	localIndex := binary.LittleEndian.Uint16(vm.readCode())

	local, err := vm.getLocal(localIndex)
	if err != nil {
		return err
	}

	return vm.push(local)
}

func (vm *vm) opPushGlobal() error {
	globalIndex := binary.LittleEndian.Uint16(vm.readCode())

	global, err := vm.getGlobal(globalIndex)
	if err != nil {
		return err
	}

	return vm.push(global)
}

func (vm *vm) opPop() error {
	_, err := vm.pop()
	return err
}

func (vm *vm) opPopLocal() error {
	localIndex := binary.LittleEndian.Uint16(vm.readCode())

	local, err := vm.pop()
	if err != nil {
		return err
	}

	err = vm.setLocal(localIndex, local)
	if err != nil {
		return err
	}

	return nil
}

func (vm *vm) opPopGlobal() error {
	globalIndex := binary.LittleEndian.Uint16(vm.readCode())

	global, err := vm.pop()
	if err != nil {
		return err
	}

	err = vm.setGlobal(globalIndex, global)
	if err != nil {
		return err
	}

	return nil
}

func (vm *vm) opInfixOperation(operation tokens.TokenType) error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalInfixOperation(left, right, operation)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opPrefixOperation(operation tokens.TokenType) error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalPrefixOperation(operand, operation)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opBranch() {
	offset := binary.LittleEndian.Uint16(vm.readCode())
	vm.movePC(int(offset))
}

func (vm *vm) opBranchIfFalse() error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	boolOperand, ok := operand.(objects.Boolean)
	if !ok {
		return errors.ErrInvalidOperand{Operation: instruction.BranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		vm.opBranch()
	}

	return nil
}

func (vm *vm) opNewArray() error {
	size := binary.LittleEndian.Uint16(vm.readCode())

	array, err := vm.popMany(int(size))
	if err != nil {
		return err
	}

	return vm.push(objects.Array(array))
}

func (vm *vm) opCollectionAccess() error {
	indexObject, err := vm.pop()
	if err != nil {
		return err
	}

	collectionObject, err := vm.pop()
	if err != nil {
		return err
	}

	var result objects.Object
	switch collection := collectionObject.(type) {
	case objects.Array:
		index, ok := indexObject.(objects.Integer)
		if !ok {
			return errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: indexObject}
		}

		if index < 0 {
			index = objects.Integer(len(collection)) + index
		}

		if index < 0 || index >= objects.Integer(len(collection)) {
			return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(collection)}
		}

		result = collection[index]
	default:
		return errors.ErrInvalidOperand{Operation: "index", LeftOperand: collection, RightOperand: indexObject}
	}

	return vm.push(result)
}

func (vm *vm) opNewFunction() error {
	templateIndex := binary.LittleEndian.Uint16(vm.readCode())

	template := vm.program.FunctionsPool[templateIndex]

	captures := make([]*objects.Object, len(template.CapturesInfo))
	for i, captureInfo := range template.CapturesInfo {
		capturedVariableScope := vm.scopeStack[len(vm.scopeStack)-1]
		for level := captureInfo.Level; level > 0; level-- {
			capturedVariableScope = capturedVariableScope.parent
		}

		captures[i] = &capturedVariableScope.variables[captureInfo.Index]
	}

	function := bytecode.Function{
		TemplateIndex: int(templateIndex),
		Captures:      captures,
	}

	return vm.push(function)
}

func (vm *vm) opCall() error {
	functionObject, err := vm.pop()
	if err != nil {
		return err
	}

	function, ok := functionObject.(bytecode.Function)
	if !ok {
		return errors.ErrInvalidOperand{Operation: "call", LeftOperand: functionObject}
	}

	functionTemplate := vm.program.FunctionsPool[function.TemplateIndex]

	arguments, err := vm.popMany(len(functionTemplate.Parameters))
	if err != nil {
		return err
	}

	functionScope := &scope{
		code:      functionTemplate.Code,
		variables: slices.Clone(arguments),
		parent:    vm.scopeStack[len(vm.scopeStack)-1],
	}

	vm.scopeStack = append(vm.scopeStack, functionScope)

	return nil
}

func (vm *vm) opReturn() error {
	vm.scopeStack = vm.scopeStack[:len(vm.scopeStack)-1]

	return nil
}
