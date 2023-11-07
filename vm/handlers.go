package vm

import (
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/exec/errors"
	"banek/exec/objects"
	"banek/exec/operations"
	"sync"
	"unsafe"
)

func (vm *vm) opPushDup(_ *scope) error {
	value, err := vm.peek()
	if err != nil {
		return err
	}

	return vm.push(value)
}

func (vm *vm) opPushConst(scope *scope) error {
	constIndex := scope.readOperand(2)

	constant, err := vm.getConst(constIndex)
	if err != nil {
		return err
	}

	return vm.push(constant.Clone())
}

func (vm *vm) opPushLocal(scope *scope) error {
	localIndex := scope.readOperand(1)

	local, err := scope.getLocal(localIndex)
	if err != nil {
		return err
	}

	return vm.push(local)
}

func (vm *vm) opPushGlobal(scope *scope) error {
	globalIndex := scope.readOperand(1)

	global, err := vm.getGlobal(globalIndex)
	if err != nil {
		return err
	}

	return vm.push(global)
}

func (vm *vm) opPushBuiltin(scope *scope) error {
	index := scope.readOperand(1)
	if index >= len(objects.Builtins) {
		return nil // TODO: return error
	}

	builtin := objects.Builtins[index]

	return vm.push(builtin)
}

func (vm *vm) opPushCaptured(scope *scope) error {
	capturedIndex := scope.readOperand(1)

	captured, err := scope.getCaptured(capturedIndex)
	if err != nil {
		return err
	}

	return vm.push(captured)
}

func (vm *vm) opPushCollElem(_ *scope) error {
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

func (vm *vm) opPop(_ *scope) error {
	_, err := vm.pop()
	return err
}

func (vm *vm) opPopLocal(scope *scope) error {
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

func (vm *vm) opPopGlobal(scope *scope) error {
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

func (vm *vm) opPopCaptured(scope *scope) error {
	capturedIndex := scope.readOperand(1)

	captured, err := vm.pop()
	if err != nil {
		return err
	}

	return scope.setCaptured(capturedIndex, captured)
}

func (vm *vm) opPopCollElem(_ *scope) error {
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

func (vm *vm) opBinaryOp(scope *scope) error {
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

func (vm *vm) opUnaryOp(scope *scope) error {
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

func (vm *vm) opBranch(scope *scope) error {
	offset := scope.readOperand(2)

	scope.movePC(offset)

	return nil
}

func (vm *vm) opBranchIfFalse(scope *scope) error {
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

func (vm *vm) opNewArray(scope *scope) error {
	size := scope.readOperand(2)

	arr := make(objects.Array, size)

	err := vm.popMany(arr)
	if err != nil {
		return err
	}

	return vm.push(arr)
}

func (vm *vm) opNewFunc(scope *scope) error {
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

var objectArrayPools = [...]sync.Pool{
	{New: func() any { return (*objects.Object)(nil) }},
	{New: func() any { return &(new([1]objects.Object)[0]) }},
	{New: func() any { return &(new([2]objects.Object)[0]) }},
	{New: func() any { return &(new([3]objects.Object)[0]) }},
	{New: func() any { return &(new([4]objects.Object)[0]) }},
}

func getObjectArray(size int) []objects.Object {
	arr := objectArrayPools[size].Get().(*objects.Object)

	return unsafe.Slice(arr, size)
}

func returnObjectArray(arr []objects.Object) {
	objectArrayPools[len(arr)].Put(unsafe.SliceData(arr))
}

func (vm *vm) opCall(scope *scope) error {
	numArgs := scope.readOperand(1)

	funcObject, err := vm.pop()
	if err != nil {
		return err
	}

	args := getObjectArray(numArgs)
	err = vm.popMany(args)
	if err != nil {
		return err
	}

	switch function := funcObject.(type) {
	case *bytecode.Func:
		funcTemplate := vm.program.FuncsPool[function.TemplateIndex]

		if len(args) > len(funcTemplate.Params) {
			args = args[:len(funcTemplate.Params)]
		}

		var locals []objects.Object
		if funcTemplate.NumLocals > len(args) {
			locals = getObjectArray(funcTemplate.NumLocals)
			copy(locals, args)
			for i := len(args); i < len(locals); i++ {
				locals[i] = objects.Undefined{}
			}
		} else {
			locals = args
		}

		vm.pushScope(funcTemplate.Code, locals, function)

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

func (vm *vm) opReturn(_ *scope) error {
	vm.popScope()

	return nil
}
