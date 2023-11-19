package vm

import (
	"banek/bytecode"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/ops"
)

func (vm *vm) opPushDup() {
	vm.push(vm.peek())
}

func (vm *vm) opPushConst() {
	constIndex := vm.readOperand(2)
	constant := vm.program.ConstsPool[constIndex]

	vm.push(constant.Clone())
}

func (vm *vm) opPush0() {
	vm.push(objs.MakeInt(0))
}

func (vm *vm) opPush1() {
	vm.push(objs.MakeInt(1))
}

func (vm *vm) opPush2() {
	vm.push(objs.MakeInt(2))
}

func (vm *vm) opPushUndefined() {
	vm.push(objs.Obj{})
}

func (vm *vm) opPushLocal() {
	localIndex := vm.readOperand(1)
	local := vm.getLocal(localIndex)

	vm.push(local)
}

func (vm *vm) opPushLocal0() {
	local := vm.getLocal(0)

	vm.push(local)
}

func (vm *vm) opPushLocal1() {
	local := vm.getLocal(1)

	vm.push(local)
}

func (vm *vm) opPushGlobal() {
	globalIndex := vm.readOperand(1)
	global := vm.getGlobal(globalIndex)

	vm.push(global)
}

func (vm *vm) opPushBuiltin() {
	index := vm.readOperand(1)
	builtin := &builtins.Funcs[index]

	vm.push(builtin.MakeObj())
}

func (vm *vm) opPushCaptured() {
	capturedIndex := vm.readOperand(1)
	captured := *vm.activeScope.captures[capturedIndex]

	vm.push(captured)
}

func (vm *vm) opPushCollElem() {
	key := vm.pop()
	coll := vm.peek()

	value, err := ops.EvalCollGet(coll, key)
	if err != nil {
		panic(err)
	}

	vm.swap(value)
}

func (vm *vm) opPop() {
	vm.pop()
}

func (vm *vm) opPopLocal() {
	localIndex := vm.readOperand(1)
	local := vm.pop()

	vm.setLocal(localIndex, local)
}

func (vm *vm) opPopLocal0() {
	local := vm.pop()

	vm.setLocal(0, local)
}

func (vm *vm) opPopLocal1() {
	local := vm.pop()

	vm.setLocal(1, local)
}

func (vm *vm) opPopGlobal() {
	globalIndex := vm.readOperand(1)
	global := vm.pop()

	vm.setGlobal(globalIndex, global)
}

func (vm *vm) opPopCaptured() {
	capturedIndex := vm.readOperand(1)
	captured := vm.pop()

	*vm.activeScope.captures[capturedIndex] = captured
}

func (vm *vm) opPopCollElem() {
	key := vm.pop()
	coll := vm.pop()
	value := vm.pop()

	err := ops.EvalCollSet(coll, key, value)
	if err != nil {
		panic(err)
	}
}

func (vm *vm) opBinaryOp() {
	operator := ops.BinaryOperator(vm.readOperand(1))

	right := vm.pop()
	left := vm.peek()

	result, err := ops.BinaryOps[operator](left, right)
	if err != nil {
		panic(err)
	}

	vm.swap(result)
}

func (vm *vm) opUnaryOp() {
	operator := ops.UnaryOperator(vm.readOperand(1))

	operand := vm.peek()

	result, err := ops.UnaryOps[operator](operand)
	if err != nil {
		panic(err)
	}

	vm.swap(result)
}

func (vm *vm) opBranch() {
	vm.pc += vm.readOperand(2)
}

func (vm *vm) opBranchIfFalse() {
	offset := vm.readOperand(2)

	operand := vm.pop()

	if operand.Tag != objs.TypeBool {
		panic(errors.ErrNotBool{Obj: operand})
	}

	if !operand.AsBool() {
		vm.pc += offset
	}
}

func (vm *vm) opNewArray() {
	size := vm.readOperand(2)

	arr := &objs.Array{
		Slice: make([]objs.Obj, size),
	}

	vm.popMany(arr.Slice)

	vm.push(objs.MakeArray(arr))
}

func (vm *vm) opNewFunc() {
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

	vm.push(function.MakeObj())
}

func (vm *vm) opCall() {
	numArgs := vm.readOperand(1)

	funcObj := vm.pop()

	switch funcObj.Tag {
	case objs.TypeFunc:
		function := bytecode.GetFunc(funcObj)
		funcTemplate := &vm.program.FuncsPool[function.TemplateIndex]

		if numArgs > funcTemplate.NumParams {
			panic(errors.ErrTooManyArgs{Expected: funcTemplate.NumParams, Received: numArgs})
		}

		vm.activeScope.savedPC = vm.pc
		vm.code = funcTemplate.Code
		vm.pc = 0

		funcScope := vm.scopeStack.pushScope(funcTemplate, function)
		vm.locals = funcScope.vars
		vm.popMany(funcScope.vars[:numArgs])
	case objs.TypeBuiltin:
		builtin := builtins.GetBuiltin(funcObj)
		if builtin.NumArgs != -1 && numArgs != builtin.NumArgs {
			panic(errors.ErrTooManyArgs{Expected: builtin.NumArgs, Received: numArgs})
		}

		args := make([]objs.Obj, numArgs)
		vm.popMany(args)

		result, err := builtin.Func(args)
		if err != nil {
			panic(err)
		}

		vm.push(result)
	default:
		panic(errors.ErrNotCallable{Obj: funcObj})
	}
}

func (vm *vm) opReturn() {
	vm.scopeStack.popScope()
	vm.pc = vm.activeScope.savedPC
	vm.code = vm.activeScope.code
	vm.locals = vm.activeScope.vars
}

func (vm *vm) opHalt() {
	vm.halted = true
}
