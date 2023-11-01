package vm

import (
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
)

func (vm *vm) opPushConst() error {
	opInfo := instruction.PushConst.Info()

	constIndex := vm.readOperand(opInfo.Operands[0].Width)

	constant, err := vm.getConstant(constIndex)
	if err != nil {
		return err
	}

	return vm.push(constant)
}

func (vm *vm) opPushLocal() error {
	opInfo := instruction.PushLocal.Info()

	localIndex := vm.readOperand(opInfo.Operands[0].Width)

	local, err := vm.getLocal(localIndex)
	if err != nil {
		return err
	}

	return vm.push(local)
}

func (vm *vm) opPushGlobal() error {
	opInfo := instruction.PushGlobal.Info()

	globalIndex := vm.readOperand(opInfo.Operands[0].Width)

	global, err := vm.getGlobal(globalIndex)
	if err != nil {
		return err
	}

	return vm.push(global)
}

func (vm *vm) opPushBuiltin() error {
	opInfo := instruction.PushBuiltin.Info()

	index := vm.readOperand(opInfo.Operands[0].Width)
	if index >= len(objects.Builtins) {
		return nil // TODO: return error
	}

	builtin := objects.Builtins[index]

	return vm.push(builtin)
}

func (vm *vm) opPop() error {
	_, err := vm.pop()
	return err
}

func (vm *vm) opPopLocal() error {
	opInfo := instruction.PopLocal.Info()

	localIndex := vm.readOperand(opInfo.Operands[0].Width)

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
	opInfo := instruction.PopGlobal.Info()

	globalIndex := vm.readOperand(opInfo.Operands[0].Width)

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

func (vm *vm) opInfixOperation() error {
	opInfo := instruction.OperationInfix.Info()

	operation := operations.InfixOperationType(vm.readOperand(opInfo.Operands[0].Width))

	operands, err := vm.popMany(2)
	if err != nil {
		return err
	}

	left, right := operands[0], operands[1]

	result, err := operations.EvalInfixOperation(left, right, operation)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opPrefixOperation() error {
	opInfo := instruction.OperationPrefix.Info()

	operation := operations.PrefixOperationType(vm.readOperand(opInfo.Operands[0].Width))

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
	opInfo := instruction.Branch.Info()

	offset := vm.readOperand(opInfo.Operands[0].Width)

	vm.movePC(offset)
}

func (vm *vm) opBranchIfFalse() error {
	opInfo := instruction.Branch.Info()

	offset := vm.readOperand(opInfo.Operands[0].Width)

	operand, err := vm.pop()
	if err != nil {
		return err
	}

	boolOperand, ok := operand.(objects.Boolean)
	if !ok {
		return errors.ErrInvalidOperand{Operation: instruction.BranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		vm.movePC(offset)
	}

	return nil
}

func (vm *vm) opNewArray() error {
	opInfo := instruction.NewArray.Info()

	size := vm.readOperand(opInfo.Operands[0].Width)

	array, err := vm.popMany(size)
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
			return errors.ErrInvalidOperand{Operation: "Index", LeftOperand: collection, RightOperand: indexObject}
		}

		if index < 0 {
			index = objects.Integer(len(collection)) + index
		}

		if index < 0 || index >= objects.Integer(len(collection)) {
			return objects.ErrIndexOutOfBounds{Index: int(index), Size: len(collection)}
		}

		result = collection[index]
	default:
		return errors.ErrInvalidOperand{Operation: "Index", LeftOperand: collection, RightOperand: indexObject}
	}

	return vm.push(result)
}

func (vm *vm) opNewFunction() error {
	opInfo := instruction.NewFunction.Info()

	templateIndex := vm.readOperand(opInfo.Operands[0].Width)

	template := vm.program.FunctionsPool[templateIndex]

	captures := make([]*objects.Object, len(template.CapturesInfo))
	for i, captureInfo := range template.CapturesInfo {
		capturedVariableScope := vm.currentScope
		for j := 0; j < captureInfo.Level; j++ {
			capturedVariableScope = capturedVariableScope.parent
		}

		captures[i] = &capturedVariableScope.variables[captureInfo.Index]
	}

	function := &bytecode.Function{
		TemplateIndex: templateIndex,
		Captures:      captures,
	}

	return vm.push(function)
}

func (vm *vm) opCall() error {
	opInfo := instruction.Call.Info()

	functionObject, err := vm.pop()
	if err != nil {
		return err
	}

	argumentsCount := vm.readOperand(opInfo.Operands[0].Width)

	switch function := functionObject.(type) {
	case *bytecode.Function:
		functionTemplate := vm.program.FunctionsPool[function.TemplateIndex]

		arguments, err := vm.popMany(argumentsCount)
		if err != nil {
			return err
		}

		if len(arguments) > len(functionTemplate.Parameters) {
			arguments = arguments[:len(functionTemplate.Parameters)]
		} else if len(arguments) < len(functionTemplate.Parameters) {
			newArguments := make([]objects.Object, len(functionTemplate.Parameters))

			copy(newArguments, arguments)
			for i := len(arguments); i < len(functionTemplate.Parameters); i++ {
				newArguments[i] = objects.Undefined
			}
		}

		functionScope := scopePool.Get().(*scope)
		functionScope.code = functionTemplate.Code
		functionScope.pc = 0
		functionScope.variables = arguments
		functionScope.parent = vm.currentScope

		vm.currentScope = functionScope

		return nil
	case objects.BuiltinFunction:
		arguments, err := vm.popMany(argumentsCount)
		if err != nil {
			return err
		}

		result, err := function.Function(arguments...)
		if err != nil {
			return err
		}

		return vm.push(result)
	default:
		return errors.ErrInvalidOperand{Operation: "call", LeftOperand: functionObject}
	}
}

func (vm *vm) opReturn() error {
	removedScope := vm.currentScope

	vm.currentScope = vm.currentScope.parent

	scopePool.Put(removedScope)

	return nil
}
