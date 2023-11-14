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

func (vm *vm) opPushDup(_ *scope) error {
	return vm.push(vm.peek())
}

func (vm *vm) opPushConst(scope *scope) error {
	constIndex := scope.readOperand(2)
	constant := vm.program.ConstsPool[constIndex]

	return vm.push(constant.Clone())
}

func (vm *vm) opPushLocal(scope *scope) error {
	localIndex := scope.readOperand(1)
	local := scope.getLocal(localIndex)

	return vm.push(local)
}

func (vm *vm) opPushGlobal(scope *scope) error {
	globalIndex := scope.readOperand(1)
	global := vm.getGlobal(globalIndex)

	return vm.push(global)
}

func (vm *vm) opPushBuiltin(scope *scope) error {
	index := scope.readOperand(1)
	builtin := builtins.Funcs[index]

	return vm.push(builtin)
}

func (vm *vm) opPushCaptured(scope *scope) error {
	capturedIndex := scope.readOperand(1)
	captured := scope.getCaptured(capturedIndex)

	return vm.push(captured)
}

func (vm *vm) opPushCollElem(_ *scope) error {
	key := vm.popOne()
	coll := vm.popOne()

	value, err := ops.EvalCollGet(coll, key)
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPop(_ *scope) error {
	vm.popOne()
	return nil
}

func (vm *vm) opPopLocal(scope *scope) error {
	localIndex := scope.readOperand(1)
	local := vm.popOne()

	scope.setLocal(localIndex, local)

	return nil
}

func (vm *vm) opPopGlobal(scope *scope) error {
	globalIndex := scope.readOperand(1)
	global := vm.popOne()

	vm.setGlobal(globalIndex, global)

	return nil
}

func (vm *vm) opPopCaptured(scope *scope) error {
	capturedIndex := scope.readOperand(1)
	captured := vm.popOne()

	scope.setCaptured(capturedIndex, captured)

	return nil
}

func (vm *vm) opPopCollElem(_ *scope) error {
	key := vm.popOne()
	coll := vm.popOne()
	value := vm.popOne()

	return ops.EvalCollSet(coll, key, value)
}

func (vm *vm) opBinaryOp(scope *scope) error {
	operator := ops.BinaryOperator(scope.readOperand(1))

	right := vm.popOne()
	left := vm.popOne()

	result, err := ops.BinaryOps[operator](left, right)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opUnaryOp(scope *scope) error {
	operator := ops.UnaryOperator(scope.readOperand(1))

	operand := vm.popOne()

	result, err := ops.UnaryOps[operator](operand)
	if err != nil {
		return err
	}

	return vm.push(result)
}

func (vm *vm) opBranch(scope *scope) error {
	offset := scope.readOperand(2)

	scope.movePC(offset)

	return nil
}

func (vm *vm) opBranchIfFalse(scope *scope) error {
	offset := scope.readOperand(2)

	operand := vm.popOne()

	boolOperand, ok := operand.(objs.Bool)
	if !ok {
		return errors.ErrInvalidOp{Operator: instrs.OpBranchIfFalse.String(), LeftOperand: boolOperand}
	}

	if !boolOperand {
		scope.movePC(offset)
	}

	return nil
}

func (vm *vm) opNewArray(scope *scope) error {
	size := scope.readOperand(2)

	arr := &objs.Array{
		Slice: make([]types.Obj, size),
	}

	vm.popMany(arr.Slice)

	return vm.push(arr)
}

func (vm *vm) opNewFunc(scope *scope) error {
	templateIndex := scope.readOperand(2)

	template := vm.program.FuncsPool[templateIndex]

	captures := make([]*types.Obj, len(template.Captures))
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

func (vm *vm) opCall(scope *scope) error {
	numArgs := scope.readOperand(1)

	funcObject := vm.popOne()

	switch function := funcObject.(type) {
	case *bytecode.Func:
		funcTemplate := vm.program.FuncsPool[function.TemplateIndex]

		if numArgs > len(funcTemplate.Params) {
			return errors.ErrTooManyArgs{Expected: len(funcTemplate.Params), Received: numArgs}
		}

		funcScope := vm.pushScope(funcTemplate.Code, funcTemplate.NumLocals, function, funcTemplate)
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

func (vm *vm) opReturn(_ *scope) error {
	vm.popScope()

	return nil
}
