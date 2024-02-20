package emulator

import (
	"banek/bytecode/instrs"
	"banek/emulator/function"
	"banek/runtime"
	"banek/runtime/binaryops"
	"banek/runtime/builtins"
	"banek/runtime/objs"
	"banek/runtime/unaryops"
	"errors"
	"unsafe"
)

func (e *emulator) handleDup()  { e.opStack.Dup() }
func (e *emulator) handleDup2() { e.opStack.Dup2() }
func (e *emulator) handleSwap() { e.opStack.Swap() }

func (e *emulator) handleJump() {
	offset := e.readOperand(instrs.OpJump, 0)
	e.movePC(offset)
}

func (e *emulator) handleBranchFalse() {
	offset := e.readOperand(instrs.OpBranchFalse, 0)

	if !e.opStack.Pop().Truthy() {
		e.movePC(offset)
	}
}

func (e *emulator) handleBranchTrue() {
	offset := e.readOperand(instrs.OpBranchTrue, 0)

	if e.opStack.Pop().Truthy() {
		e.movePC(offset)
	}
}

func (e *emulator) handlePushBuiltin() {
	index := e.readOperand(instrs.OpPushBuiltin, 0)
	e.opStack.Push(objs.Obj{Type: objs.Builtin, Int: index})
}

func (e *emulator) handlePushGlobal() {
	index := e.readOperand(instrs.OpPushGlobal, 0)
	e.opStack.Push(e.opStack.ReadVar(e.callStack.GlobalFrame().BP, index))
}

func (e *emulator) handlePushLocal() {
	index := e.readOperand(instrs.OpPushLocal, 0)
	e.opStack.Push(e.opStack.ReadVar(e.callStack.ActiveFrame().BP, index))
}

func (e *emulator) handlePushLocal0() {
	e.opStack.Push(e.opStack.ReadVar(e.callStack.ActiveFrame().BP, 0))
}

func (e *emulator) handlePushLocal1() {
	e.opStack.Push(e.opStack.ReadVar(e.callStack.ActiveFrame().BP, 1))
}

func (e *emulator) handlePushLocal2() {
	e.opStack.Push(e.opStack.ReadVar(e.callStack.ActiveFrame().BP, 2))
}

func (e *emulator) handlePush0()  { e.opStack.Push(objs.MakeInt(0)) }
func (e *emulator) handlePush1()  { e.opStack.Push(objs.MakeInt(1)) }
func (e *emulator) handlePush2()  { e.opStack.Push(objs.MakeInt(2)) }
func (e *emulator) handlePush3()  { e.opStack.Push(objs.MakeInt(3)) }
func (e *emulator) handlePushN1() { e.opStack.Push(objs.MakeInt(-1)) }

func (e *emulator) handlePushInt() {
	value := e.readOperand(instrs.OpPushInt, 0)

	e.opStack.Push(objs.MakeInt(value))
}

func (e *emulator) handlePushStr() {
	index := e.readOperand(instrs.OpPushStr, 0)
	str := e.program.StringPool[index]

	e.opStack.Push(objs.MakeString(str))
}

func (e *emulator) handlePushTrue()  { e.opStack.Push(objs.MakeBool(true)) }
func (e *emulator) handlePushFalse() { e.opStack.Push(objs.MakeBool(false)) }
func (e *emulator) handlePushUndef() { e.opStack.Push(objs.Obj{}) }

func (e *emulator) handlePop() { e.opStack.Pop() }

func (e *emulator) handlePopGlobal() {
	index := e.readOperand(instrs.OpPopGlobal, 0)
	e.opStack.WriteVar(e.callStack.GlobalFrame().BP, index, e.opStack.Pop())
}

func (e *emulator) handlePopLocal() {
	index := e.readOperand(instrs.OpPopLocal, 0)
	e.opStack.WriteVar(e.callStack.ActiveFrame().BP, index, e.opStack.Pop())
}

func (e *emulator) handlePopLocal0() {
	e.opStack.WriteVar(e.callStack.ActiveFrame().BP, 0, e.opStack.Pop())
}

func (e *emulator) handlePopLocal1() {
	e.opStack.WriteVar(e.callStack.ActiveFrame().BP, 1, e.opStack.Pop())
}

func (e *emulator) handlePopLocal2() {
	e.opStack.WriteVar(e.callStack.ActiveFrame().BP, 2, e.opStack.Pop())
}

func (e *emulator) handleMakeArray() {
	size := e.readOperand(instrs.OpMakeArray, 0)

	arr := make([]objs.Obj, size)
	e.opStack.PopMany(arr)

	e.opStack.Push(objs.MakeArray(arr))
}

func (e *emulator) handleNewArray() {
	sizeObj := e.opStack.Pop()
	if sizeObj.Type != objs.Int {
		// TODO: handle non-integer size
	}

	arr := make([]objs.Obj, sizeObj.Int)

	e.opStack.Push(objs.MakeArray(arr))
}

func (e *emulator) handleBinaryAdd() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.AddOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleBinarySub() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.SubOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleBinaryMul() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.MulOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleBinaryDiv() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.DivOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleBinaryMod() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.ModOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleBinaryEq() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	e.opStack.Push(objs.MakeBool(left.Equals(right)))
}

func (e *emulator) handleBinaryNeq() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	e.opStack.Push(objs.MakeBool(!left.Equals(right)))
}

func (e *emulator) handleCompLt() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := left.Compare(right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(objs.MakeBool(result < 0))
}

func (e *emulator) handleCompLe() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := left.Compare(right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(objs.MakeBool(result <= 0))
}

func (e *emulator) handleCompGt() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := left.Compare(right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(objs.MakeBool(result > 0))
}

func (e *emulator) handleCompGe() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := left.Compare(right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(objs.MakeBool(result >= 0))
}

func (e *emulator) handleUnaryNeg() {
	operand := e.opStack.Pop()

	result, err := unaryops.NegOperator.Eval(operand)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleUnaryNot() {
	operand := e.opStack.Pop()

	result, err := unaryops.NotOperator.Eval(operand)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleMakeFunc() {
	templateIndex := e.readOperand(instrs.OpMakeFunc, 0)
	template := &e.program.FuncPool[templateIndex]

	captures := make([]*objs.Obj, len(template.Captures))
	for i, index := range template.Captures {
		// TODO: handle captured globals
		_, _ = i, index
	}

	fn := &function.Obj{
		Index:    templateIndex,
		Captures: captures,
	}

	e.opStack.Push(objs.Obj{Type: objs.Func, Ptr: unsafe.Pointer(fn)})
}

func (e *emulator) handleCall() {
	numArgs := e.readOperand(instrs.OpCall, 0)

	funcObj := e.opStack.Pop()

	switch funcObj.Type {
	case objs.Func:
		fn := (*function.Obj)(funcObj.Ptr)
		template := &e.program.FuncPool[fn.Index]

		if numArgs > template.NumParams {
			panic(runtime.TooManyArgsError{Expected: template.NumParams, Actual: numArgs})
		}

		e.callStack.Push(template.Addr, e.opStack.SP()-numArgs, fn)
		e.opStack.Grow(template.NumLocals - numArgs)
	case objs.Builtin:
		builtinIndex := funcObj.Int
		builtin := &builtins.Funcs[builtinIndex]

		if builtin.NumArgs != -1 && numArgs != builtin.NumArgs {
			panic(runtime.TooManyArgsError{Expected: builtin.NumArgs, Actual: numArgs})
		}

		args := make([]objs.Obj, numArgs)
		e.opStack.PopMany(args)

		res, err := builtin.Func(args)
		if err != nil {
			err = errors.Join(builtins.CallError{BuiltinName: builtin.Name}, err)
			panic(err)
		}

		e.opStack.Push(res)
	default:
		panic(runtime.NotCallableError{Func: funcObj})
	}
}

func (e *emulator) handleReturn() {
	e.opStack.PatchReturn(e.callStack.ActiveFrame().BP)
	e.callStack.Pop()
}

func (_ *emulator) handleHalt() {
	panic(struct{}{})
}

func (e *emulator) handlePushCollElem() {
	key := e.opStack.Pop()
	coll := e.opStack.Pop()

	elem, err := coll.GetField(key)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(elem)
}

func (e *emulator) handlePopCollElem() {
	value := e.opStack.Pop()
	key := e.opStack.Pop()
	coll := e.opStack.Pop()

	err := coll.SetField(key, value)
	if err != nil {
		panic(err)
	}
}
