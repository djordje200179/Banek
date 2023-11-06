package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
)

func (vm *vm) opPushDup(_ *Scope) error {
	value, err := vm.peek()
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPushConst(scope *Scope) error {
	constIndex := scope.readOperand(2)

	constant, err := vm.getConst(constIndex)
	if err != nil {
		return err
	}

	return vm.push(constant.Clone())
}

func (vm *vm) opPushLocal(scope *Scope) error {
	localIndex := scope.readOperand(1)

	local, err := scope.getLocal(localIndex)
	if err != nil {
		return err
	}

	return vm.push(local)
}

func (vm *vm) opPushGlobal(scope *Scope) error {
	globalIndex := scope.readOperand(1)

	global, err := vm.getGlobal(globalIndex)
	if err != nil {
		return err
	}

	return vm.push(global)
}

func (vm *vm) opPushBuiltin(scope *Scope) error {
	index := scope.readOperand(1)
	if index >= len(objects.Builtins) {
		return nil // TODO: return error
	}

	builtin := objects.Builtins[index]

	return vm.push(builtin)
}

func (vm *vm) opPushCaptured(scope *Scope) error {
	capturedIndex := scope.readOperand(1)

	captured, err := scope.getCaptured(capturedIndex)
	if err != nil {
		return err
	}

	return vm.push(captured)
}

func (vm *vm) opPushCollElem(_ *Scope) error {
	key, err := vm.pop()
	if err != nil {
		return err
	}

	coll, err := vm.pop()
	if err != nil {
		return err
	}

	value, err := operations.EvalCollGet(coll, key)
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPop(_ *Scope) error {
	_, err := vm.pop()
	return err
}

func (vm *vm) opPopLocal(scope *Scope) error {
	localIndex := scope.readOperand(1)

	local, err := vm.pop()
	if err != nil {
		return err
	}

	err = scope.setLocal(localIndex, local)
	if err != nil {
		return err
	}

	return nil
}

func (vm *vm) opPopGlobal(scope *Scope) error {
	globalIndex := scope.readOperand(1)

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

func (vm *vm) opPopCaptured(scope *Scope) error {
	capturedIndex := scope.readOperand(1)

	captured, err := vm.pop()
	if err != nil {
		return err
	}

	return scope.setCaptured(capturedIndex, captured)
}

func (vm *vm) opPopCollElem(_ *Scope) error {
	key, err := vm.pop()
	if err != nil {
		return err
	}

	coll, err := vm.pop()
	if err != nil {
		return err
	}

	value, err := vm.pop()
	if err != nil {
		return err
	}

	return operations.EvalCollSet(coll, key, value)
}

func (vm *vm) opBinaryOp(scope *Scope) error {
	operator := operations.BinaryOperator(scope.readOperand(1))

	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalBinary(left, right, operator)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opUnaryOp(scope *Scope) error {
	operator := operations.UnaryOperator(scope.readOperand(1))

	operand, err := vm.pop()
	if err != nil {
		return err
	}

	result, err := operations.EvalUnary(operand, operator)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opBranch(scope *Scope) error {
	offset := scope.readOperand(2)

	scope.movePC(offset)

	return nil
}

func (vm *vm) opBranchIfFalse(scope *Scope) error {
	offset := scope.readOperand(2)

	operand, err := vm.pop()
	if err != nil {
		return err
	}

	boolOperand, ok := operand.(objects.Boolean)
	if !ok {
		return errors.ErrInvalidOp{Operator: instructions.OpBranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		scope.movePC(offset)
	}

	return nil
}

func (vm *vm) opNewArray(scope *Scope) error {
	size := scope.readOperand(2)

	arr := make(objects.Array, size)

	err := vm.popMany(arr)
	if err != nil {
		return err
	}

	return vm.push(arr)
}

func (vm *vm) opNewFunc(scope *Scope) error {
	templateIndex := scope.readOperand(2)

	template := vm.program.FuncsPool[templateIndex]

	captures := make([]*objects.Object, len(template.Captures))
	for i, captureInfo := range template.Captures {
		capturedVariableScope := vm.currScope
		for j := 0; j < captureInfo.Level; j++ {
			capturedVariableScope = capturedVariableScope.parent
		}

		captures[i] = &capturedVariableScope.vars[captureInfo.Index]
	}

	function := &bytecode.Func{
		TemplateIndex: templateIndex,
		Captures:      captures,
	}

	return vm.push(function)
}

func (vm *vm) opCall(scope *Scope) error {
	numArgs := scope.readOperand(1)

	funcObject, err := vm.pop()
	if err != nil {
		return err
	}

	args := make([]objects.Object, numArgs)
	err = vm.popMany(args)
	if err != nil {
		return err
	}

	switch function := funcObject.(type) {
	case *bytecode.Func:
		vm.pushScope(function, args)

		return nil
	case objects.BuiltinFunc:
		result, err := function.Func(args...)
		if err != nil {
			return err
		}

		return vm.push(result)
	default:
		return errors.ErrInvalidOp{Operator: "call", LeftOperand: funcObject}
	}
}

func (vm *vm) opReturn(_ *Scope) error {
	vm.popScope()

	return nil
}
