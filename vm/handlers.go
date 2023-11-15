package vm

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/ops"
	"banek/runtime/types"
)

func (vm *vm) opPushDup() error {
	return vm.push(vm.peek())
}

func (vm *vm) opPushConst() error {
	constIndex := vm.readOperand(2)
	constant := vm.program.ConstsPool[constIndex]

	return vm.push(constant.Clone())
}

func (vm *vm) opPush0() error {
	return vm.push(objs.Int(0))
}

func (vm *vm) opPush1() error {
	return vm.push(objs.Int(1))
}

func (vm *vm) opPush2() error {
	return vm.push(objs.Int(2))
}

func (vm *vm) opPushUndefined() error {
	return vm.push(objs.Undefined{})
}

func (vm *vm) opPushLocal() error {
	localIndex := vm.readOperand(1)
	local := vm.getLocal(localIndex)

	return vm.push(local)
}

func (vm *vm) opPushLocal0() error {
	local := vm.getLocal(0)

	return vm.push(local)
}

func (vm *vm) opPushLocal1() error {
	local := vm.getLocal(1)

	return vm.push(local)
}

func (vm *vm) opPushGlobal() error {
	globalIndex := vm.readOperand(1)
	global := vm.getGlobal(globalIndex)

	return vm.push(global)
}

func (vm *vm) opPushBuiltin() error {
	index := vm.readOperand(1)
	builtin := builtins.Funcs[index]

	return vm.push(builtin)
}

func (vm *vm) opPushCaptured() error {
	capturedIndex := vm.readOperand(1)
	captured := vm.getCaptured(capturedIndex)

	return vm.push(captured)
}

func (vm *vm) opPushCollElem() error {
	key := vm.pop()
	coll := vm.pop()

	value, err := ops.EvalCollGet(coll, key)
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPop() error {
	vm.pop()
	return nil
}

func (vm *vm) opPopLocal() error {
	localIndex := vm.readOperand(1)
	local := vm.pop()

	vm.setLocal(localIndex, local)

	return nil
}

func (vm *vm) opPopLocal0() error {
	local := vm.pop()

	vm.setLocal(0, local)

	return nil
}

func (vm *vm) opPopLocal1() error {
	local := vm.pop()

	vm.setLocal(1, local)

	return nil
}

func (vm *vm) opPopGlobal() error {
	globalIndex := vm.readOperand(1)
	global := vm.pop()

	vm.setGlobal(globalIndex, global)

	return nil
}

func (vm *vm) opPopCaptured() error {
	capturedIndex := vm.readOperand(1)
	captured := vm.pop()

	vm.setCaptured(capturedIndex, captured)

	return nil
}

func (vm *vm) opPopCollElem() error {
	key := vm.pop()
	coll := vm.pop()
	value := vm.pop()

	return ops.EvalCollSet(coll, key, value)
}

func (vm *vm) opBinaryOp() error {
	operator := ops.BinaryOperator(vm.readOperand(1))

	right := vm.pop()
	left := vm.pop()

	result, err := ops.BinaryOps[operator](left, right)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opUnaryOp() error {
	operator := ops.UnaryOperator(vm.readOperand(1))

	operand := vm.pop()

	result, err := ops.UnaryOps[operator](operand)
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

	operand := vm.pop()

	boolOperand, ok := operand.(objs.Bool)
	if !ok {
		return errors.ErrInvalidOp{Operator: instrs.OpBranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		vm.movePC(offset)
	}

	return nil
}

func (vm *vm) opNewArray() error {
	size := vm.readOperand(2)

	arr := &objs.Array{
		Slice: make([]types.Obj, size),
	}

	vm.popMany(arr.Slice)

	return vm.push(arr)
}

func (vm *vm) opNewFunc() error {
	templateIndex := vm.readOperand(2)

	template := vm.program.FuncsPool[templateIndex]

	captures := make([]*types.Obj, len(template.Captures))
	for i, captureInfo := range template.Captures {
		capturedVariableScope := vm.scope
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

	funcObject := vm.pop()

	switch function := funcObject.(type) {
	case *bytecode.Func:
		funcTemplate := vm.program.FuncsPool[function.TemplateIndex]

		if numArgs > len(funcTemplate.Params) {
			return errors.ErrTooManyArgs{Expected: len(funcTemplate.Params), Received: numArgs}
		}

		funcScope := vm.pushScope(funcTemplate.Code, funcTemplate.NumLocals, function)
		vm.popMany(funcScope.vars[:numArgs])

		return nil
	case builtins.BuiltinFunc:
		if function.NumArgs != -1 && numArgs != function.NumArgs {
			return errors.ErrTooManyArgs{Expected: function.NumArgs, Received: numArgs}
		}

		args := make([]types.Obj, numArgs)
		vm.popMany(args)

		result, err := function.Func(args)
		if err != nil {
			return err
		}

		return vm.push(result)
	default:
		return errors.ErrNotCallable{Obj: funcObject}
	}
}

func (vm *vm) opReturn() error {
	funcTemplateIndex := vm.scope.function.TemplateIndex
	funcTemplate := vm.program.FuncsPool[funcTemplateIndex]
	canFreeVars := !funcTemplate.IsCaptured

	vm.popScope(canFreeVars)

	return nil
}

func (vm *vm) opHalt() error {
	vm.halted = true

	return nil
}
