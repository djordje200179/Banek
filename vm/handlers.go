package vm

import (
	"banek/bytecode"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/ops"
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
	return vm.push(objs.MakeInt(0))
}

func (vm *vm) opPush1() error {
	return vm.push(objs.MakeInt(1))
}

func (vm *vm) opPush2() error {
	return vm.push(objs.MakeInt(2))
}

func (vm *vm) opPushUndefined() error {
	return vm.push(objs.Obj{})
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
	builtin := &builtins.Funcs[index]

	return vm.push(builtin.MakeObj())
}

func (vm *vm) opPushCaptured() error {
	capturedIndex := vm.readOperand(1)
	captured := *vm.activeScope.captures[capturedIndex]

	return vm.push(captured)
}

func (vm *vm) opPushCollElem() error {
	key := vm.pop()
	coll := vm.peek()

	value, err := ops.EvalCollGet(coll, key)
	if err != nil {
		return err
	}

	vm.swap(value)

	return nil
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

	*vm.activeScope.captures[capturedIndex] = captured

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
	left := vm.peek()

	result, err := ops.BinaryOps[operator](left, right)
	if err != nil {
		return err
	}

	vm.swap(result)

	return nil
}

func (vm *vm) opUnaryOp() error {
	operator := ops.UnaryOperator(vm.readOperand(1))

	operand := vm.peek()

	result, err := ops.UnaryOps[operator](operand)
	if err != nil {
		return err
	}

	vm.swap(result)

	return nil
}

func (vm *vm) opBranch() error {
	vm.pc += vm.readOperand(2)

	return nil
}

func (vm *vm) opBranchIfFalse() error {
	offset := vm.readOperand(2)

	operand := vm.pop()

	if operand.Tag != objs.TypeBool {
		return errors.ErrNotBool{Obj: operand}
	}

	if !operand.AsBool() {
		vm.pc += offset
	}

	return nil
}

func (vm *vm) opNewArray() error {
	size := vm.readOperand(2)

	arr := &objs.Array{
		Slice: make([]objs.Obj, size),
	}

	vm.popMany(arr.Slice)

	return vm.push(objs.MakeArray(arr))
}

func (vm *vm) opNewFunc() error {
	templateIndex := vm.readOperand(2)

	template := vm.program.FuncsPool[templateIndex]

	captures := make([]*objs.Obj, len(template.Captures))
	for i, captureInfo := range template.Captures {
		capturedVariableScope := vm.activeScope
		for j := 0; j < captureInfo.Level; j++ {
			capturedVariableScope = capturedVariableScope.parent
		}

		captures[i] = &capturedVariableScope.vars[captureInfo.Index]
	}

	function := &bytecode.Func{
		TemplateIndex: templateIndex,
		Captures:      captures,
	}

	return vm.push(function.MakeObj())
}

func (vm *vm) opCall() error {
	numArgs := vm.readOperand(1)

	funcObj := vm.pop()

	switch funcObj.Tag {
	case objs.TypeFunc:
		function := bytecode.GetFunc(funcObj)
		funcTemplate := &vm.program.FuncsPool[function.TemplateIndex]

		if numArgs > len(funcTemplate.Params) {
			return errors.ErrTooManyArgs{Expected: len(funcTemplate.Params), Received: numArgs}
		}

		vm.activeScope.savedPC = vm.pc
		vm.code = funcTemplate.Code
		vm.pc = 0

		funcScope := vm.scopeStack.pushScope(funcTemplate, function)
		vm.locals = funcScope.vars
		vm.popMany(funcScope.vars[:numArgs])

		return nil
	case objs.TypeBuiltin:
		builtin := builtins.GetBuiltin(funcObj)
		if builtin.NumArgs != -1 && numArgs != builtin.NumArgs {
			return errors.ErrTooManyArgs{Expected: builtin.NumArgs, Received: numArgs}
		}

		args := make([]objs.Obj, numArgs)
		vm.popMany(args)

		result, err := builtin.Func(args)
		if err != nil {
			return err
		}

		return vm.push(result)
	default:
		return errors.ErrNotCallable{Obj: funcObj}
	}
}

func (vm *vm) opReturn() error {
	vm.scopeStack.popScope()
	vm.pc = vm.activeScope.savedPC
	vm.code = vm.activeScope.code
	vm.locals = vm.activeScope.vars

	return nil
}

func (vm *vm) opHalt() error {
	vm.halted = true

	return nil
}
