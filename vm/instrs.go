package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
)

func (vm *vm) opPushDup() error {
	value, err := vm.peek()
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPushConst() error {
	constIndex := vm.readOperand(2)

	constant, err := vm.getConst(constIndex)
	if err != nil {
		return err
	}

	return vm.push(constant.Clone())
}

func (vm *vm) opPushLocal() error {
	localIndex := vm.readOperand(1)

	local, err := vm.getLocal(localIndex)
	if err != nil {
		return err
	}

	return vm.push(local)
}

func (vm *vm) opPushGlobal() error {
	globalIndex := vm.readOperand(1)

	global, err := vm.getGlobal(globalIndex)
	if err != nil {
		return err
	}

	return vm.push(global)
}

func (vm *vm) opPushBuiltin() error {
	index := vm.readOperand(1)
	if index >= len(objects.Builtins) {
		return nil // TODO: return error
	}

	builtin := objects.Builtins[index]

	return vm.push(builtin)
}

func (vm *vm) opPushCaptured() error {
	capturedIndex := vm.readOperand(1)

	captured, err := vm.getCaptured(capturedIndex)
	if err != nil {
		return err
	}

	return vm.push(captured)
}

func (vm *vm) opPushCollElem() error {
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

func (vm *vm) opPop() error {
	_, err := vm.pop()
	return err
}

func (vm *vm) opPopLocal() error {
	localIndex := vm.readOperand(1)

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
	globalIndex := vm.readOperand(1)

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

func (vm *vm) opPopCaptured() error {
	capturedIndex := vm.readOperand(1)

	captured, err := vm.pop()
	if err != nil {
		return err
	}

	return vm.setCaptured(capturedIndex, captured)
}

func (vm *vm) opPopCollElem() error {
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

func (vm *vm) opBinaryOp() error {
	operator := operations.BinaryOperator(vm.readOperand(1))

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

func (vm *vm) opUnaryOp() error {
	operator := operations.UnaryOperator(vm.readOperand(1))

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

func (vm *vm) opBranch() error {
	offset := vm.readOperand(2)

	vm.movePC(offset)

	return nil
}

func (vm *vm) opBranchIfFalse() error {
	offset := vm.readOperand(2)

	operand, err := vm.pop()
	if err != nil {
		return err
	}

	boolOperand, ok := operand.(objects.Boolean)
	if !ok {
		return errors.ErrInvalidOp{Operator: instructions.OpBranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		vm.movePC(offset)
	}

	return nil
}

func (vm *vm) opNewArray() error {
	size := vm.readOperand(2)

	arr := make(objects.Array, size)

	err := vm.popMany(arr)
	if err != nil {
		return err
	}

	return vm.push(arr)
}

func (vm *vm) opNewFunc() error {
	templateIndex := vm.readOperand(2)

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

func (vm *vm) opCall() error {
	numArgs := vm.readOperand(1)

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
		funcTemplate := vm.program.FuncsPool[function.TemplateIndex]

		if len(args) < len(funcTemplate.Params) {
			oldArgs := args
			args = make([]objects.Object, len(funcTemplate.Params))

			copy(args, oldArgs)
			for i := len(oldArgs); i < len(funcTemplate.Params); i++ {
				args[i] = objects.Undefined{}
			}
		} else if len(args) > len(funcTemplate.Params) {
			args = args[:len(funcTemplate.Params)]
		}

		funcScope := scopePool.Get().(*scope)
		funcScope.code = funcTemplate.Code
		funcScope.pc = 0
		funcScope.vars = args
		funcScope.parent = vm.currScope
		funcScope.function = function

		vm.currScope = funcScope

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

func (vm *vm) opReturn() error {
	removedScope := vm.currScope

	vm.currScope = vm.currScope.parent

	scopePool.Put(removedScope)

	return nil
}
