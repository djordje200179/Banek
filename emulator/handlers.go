package emulator

import (
	"banek/bytecode/instrs"
	"banek/emulator/callstack"
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

func (e *emulator) handleBuiltin() {
	index := e.readOperand(instrs.OpBuiltin, 0)
	e.opStack.Push(objs.Obj{Type: objs.Builtin, Int: index})
}

func (e *emulator) handleLoadGlobal() {
	index := e.readOperand(instrs.OpLoadGlobal, 0)
	e.opStack.Push(e.opStack.ReadVar(0, index))
}

func (e *emulator) handleLoadLocal() {
	index := e.readOperand(instrs.OpLoadLocal, 0)
	e.opStack.Push(e.opStack.ReadVar(e.frame.BP, index))
}

func (e *emulator) handleLoadLocal0() { e.opStack.Push(e.opStack.ReadVar(e.frame.BP, 0)) }
func (e *emulator) handleLoadLocal1() { e.opStack.Push(e.opStack.ReadVar(e.frame.BP, 1)) }
func (e *emulator) handleLoadLocal2() { e.opStack.Push(e.opStack.ReadVar(e.frame.BP, 2)) }

func (e *emulator) handleConst0()  { e.opStack.Push(objs.MakeInt(0)) }
func (e *emulator) handleConst1()  { e.opStack.Push(objs.MakeInt(1)) }
func (e *emulator) handleConst2()  { e.opStack.Push(objs.MakeInt(2)) }
func (e *emulator) handleConst3()  { e.opStack.Push(objs.MakeInt(3)) }
func (e *emulator) handleConstN1() { e.opStack.Push(objs.MakeInt(-1)) }

func (e *emulator) handleConstInt() {
	value := e.readOperand(instrs.OpConstInt, 0)

	e.opStack.Push(objs.MakeInt(value))
}

func (e *emulator) handleConstStr() {
	index := e.readOperand(instrs.OpConstStr, 0)
	str := e.program.StringPool[index]

	e.opStack.Push(objs.MakeString(str))
}

func (e *emulator) handleConstTrue()  { e.opStack.Push(objs.MakeBool(true)) }
func (e *emulator) handleConstFalse() { e.opStack.Push(objs.MakeBool(false)) }
func (e *emulator) handleConstUndef() { e.opStack.Push(objs.Obj{}) }

func (e *emulator) handlePop() { e.opStack.Pop() }

func (e *emulator) handleStoreGlobal() {
	index := e.readOperand(instrs.OpStoreGlobal, 0)
	e.opStack.WriteVar(0, index, e.opStack.Pop())
}

func (e *emulator) handleStoreLocal() {
	index := e.readOperand(instrs.OpStoreLocal, 0)
	e.opStack.WriteVar(e.frame.BP, index, e.opStack.Pop())
}

func (e *emulator) handleStoreLocal0() { e.opStack.WriteVar(e.frame.BP, 0, e.opStack.Pop()) }
func (e *emulator) handleStoreLocal1() { e.opStack.WriteVar(e.frame.BP, 1, e.opStack.Pop()) }
func (e *emulator) handleStoreLocal2() { e.opStack.WriteVar(e.frame.BP, 2, e.opStack.Pop()) }

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

func (e *emulator) handleAdd() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.AddOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleSub() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.SubOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleMul() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.MulOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleDiv() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.DivOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleMod() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	result, err := binaryops.ModOperator.Eval(left, right)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleCompEq() {
	right := e.opStack.Pop()
	left := e.opStack.Pop()

	e.opStack.Push(objs.MakeBool(left.Equals(right)))
}

func (e *emulator) handleCompNeq() {
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

func (e *emulator) handleNeg() {
	operand := e.opStack.Pop()

	result, err := unaryops.NegOperator.Eval(operand)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(result)
}

func (e *emulator) handleNot() {
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

		e.callStack.Push(e.frame)
		e.frame = callstack.Frame{
			PC: template.Addr,
			BP: e.opStack.SP() - numArgs,

			Func: fn,
		}

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
	e.opStack.PatchReturn(e.frame.BP)
	e.frame = e.callStack.Pop()
}

func (_ *emulator) handleHalt() {
	panic(struct{}{})
}

func (e *emulator) handleCollGet() {
	key := e.opStack.Pop()
	coll := e.opStack.Pop()

	elem, err := coll.Get(key)
	if err != nil {
		panic(err)
	}

	e.opStack.Push(elem)
}

func (e *emulator) handleCollSet() {
	value := e.opStack.Pop()
	key := e.opStack.Pop()
	coll := e.opStack.Pop()

	err := coll.Set(key, value)
	if err != nil {
		panic(err)
	}
}
