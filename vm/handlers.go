package vm

import (
	"banek/bytecode"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
	"banek/runtime/ops"
	"banek/vm/scopes"
	"banek/vm/stack"
)

func opPushDup(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(operandStack.PopAndReserve())
}

func opPushConst(program *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	constant := program.ConstsPool[scopeStack.ReadOperand(2)]

	operandStack.Push(constant.Clone())
}

func opPush0(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(objs.MakeInt(0))
}

func opPush1(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(objs.MakeInt(1))
}

func opPush2(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(objs.MakeInt(2))
}

func opPushUndefined(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(objs.Obj{})
}

func opPushLocal(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(scopeStack.GetLocal(scopeStack.ReadOperand(1)))
}

func opPushLocal0(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(scopeStack.GetLocal(0))
}

func opPushLocal1(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(scopeStack.GetLocal(1))
}

func opPushGlobal(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(scopeStack.GetGlobal(scopeStack.ReadOperand(1)))
}

func opPushBuiltin(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	index := scopeStack.ReadOperand(1)
	builtin := &builtins.Funcs[index]

	operandStack.Push(builtin.MakeObj())
}

func opPushCaptured(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Push(scopeStack.GetCaptured(scopeStack.ReadOperand(1)))
}

func opPushCollElem(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	key := operandStack.Pop()
	coll := operandStack.PopAndReserve()

	value, err := ops.EvalCollGet(coll, key)
	if err != nil {
		panic(err)
	}

	operandStack.PushReservation(value)
}

func opPop(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	operandStack.Pop()
}

func opPopLocal(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	scopeStack.SetLocal(scopeStack.ReadOperand(1), operandStack.Pop())
}

func opPopLocal0(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	scopeStack.SetLocal(0, operandStack.Pop())
}

func opPopLocal1(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	scopeStack.SetLocal(1, operandStack.Pop())
}

func opPopGlobal(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	scopeStack.SetGlobal(scopeStack.ReadOperand(1), operandStack.Pop())
}

func opPopCaptured(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	scopeStack.SetCaptured(scopeStack.ReadOperand(1), operandStack.Pop())
}

func opPopCollElem(_ *bytecode.Executable, _ *scopes.Stack, operandStack *stack.Stack) {
	key := operandStack.Pop()
	coll := operandStack.Pop()
	value := operandStack.Pop()

	err := ops.EvalCollSet(coll, key, value)
	if err != nil {
		panic(err)
	}
}

func opBinaryOp(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operator := ops.BinaryOperator(scopeStack.ReadOperand(1))

	right := operandStack.Pop()
	left := operandStack.PopAndReserve()

	result, err := ops.BinaryOps[operator](left, right)
	if err != nil {
		panic(err)
	}

	operandStack.PushReservation(result)
}

func opUnaryOp(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	operator := ops.UnaryOperator(scopeStack.ReadOperand(1))

	operand := operandStack.PopAndReserve()

	result, err := ops.UnaryOps[operator](operand)
	if err != nil {
		panic(err)
	}

	operandStack.PushReservation(result)
}

func opBranch(_ *bytecode.Executable, scopeStack *scopes.Stack, _ *stack.Stack) {
	scopeStack.MovePC(scopeStack.ReadOperand(2))
}

func opBranchIfFalse(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	offset := scopeStack.ReadOperand(2)

	operand := operandStack.Pop()

	if operand.Tag != objs.TypeBool {
		panic(errors.ErrNotBool{Obj: operand})
	}

	if !operand.AsBool() {
		scopeStack.MovePC(offset)
	}
}

func opNewArray(_ *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	size := scopeStack.ReadOperand(2)

	arr := &objs.Array{
		Slice: make([]objs.Obj, size),
	}

	operandStack.PopMany(arr.Slice)

	operandStack.Push(objs.MakeArray(arr))
}

func opNewFunc(program *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	templateIndex := scopeStack.ReadOperand(2)

	template := program.FuncsPool[templateIndex]

	captures := make([]*objs.Obj, len(template.Captures))
	for i, captureInfo := range template.Captures {
		captures[i] = scopeStack.GetCapture(captureInfo)
	}

	function := &bytecode.Func{
		TemplateIndex: templateIndex,
		Captures:      captures,
	}

	operandStack.Push(function.MakeObj())
}

func opCall(program *bytecode.Executable, scopeStack *scopes.Stack, operandStack *stack.Stack) {
	numArgs := scopeStack.ReadOperand(1)

	funcObj := operandStack.Pop()

	switch funcObj.Tag {
	case objs.TypeFunc:
		function := bytecode.GetFunc(funcObj)
		template := &program.FuncsPool[function.TemplateIndex]

		if numArgs > template.NumParams {
			panic(errors.ErrTooManyArgs{Expected: template.NumParams, Received: numArgs})
		}

		locals := scopeStack.NewScope(function, template)

		operandStack.PopMany(locals[:numArgs])
	case objs.TypeBuiltin:
		builtin := builtins.GetBuiltin(funcObj)
		if builtin.NumArgs != -1 && numArgs != builtin.NumArgs {
			panic(errors.ErrTooManyArgs{Expected: builtin.NumArgs, Received: numArgs})
		}

		args := make([]objs.Obj, numArgs)
		operandStack.PopMany(args)

		result, err := builtin.Func(args)
		if err != nil {
			panic(err)
		}

		operandStack.Push(result)
	default:
		panic(errors.ErrNotCallable{Obj: funcObj})
	}
}

func opReturn(_ *bytecode.Executable, scopeStack *scopes.Stack, _ *stack.Stack) {
	scopeStack.RestoreScope()
}
