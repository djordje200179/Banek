package emulator

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime"
	"banek/runtime/builtins"
	"banek/runtime/primitives"
)

func (e *emulator) handleDup()  { e.stack.Dup() }
func (e *emulator) handleDup2() { e.stack.Dup2() }
func (e *emulator) handleSwap() { e.stack.Swap() }

func (e *emulator) handleJump() {
	offset := e.readOperand(instrs.OpJump, 0)
	e.movePC(offset)
}

func (e *emulator) handleBranchFalse() {
	offset := e.readOperand(instrs.OpBranchFalse, 0)

	if !e.stack.Pop().Truthy() {
		e.movePC(offset)
	}
}

func (e *emulator) handleBranchTrue() {
	offset := e.readOperand(instrs.OpBranchTrue, 0)

	if e.stack.Pop().Truthy() {
		e.movePC(offset)
	}
}

func (e *emulator) handlePushBuiltin() {
	builtin := e.readOperand(instrs.OpPushBuiltin, 0)
	e.stack.Push(&builtins.Funcs[builtin])
}

func (e *emulator) handlePushGlobal() {
	index := e.readOperand(instrs.OpPushGlobal, 0)
	e.stack.Push(e.stack.ReadVar(e.globalScope.bp, index))
}

func (e *emulator) handlePushLocal() {
	index := e.readOperand(instrs.OpPushLocal, 0)
	e.stack.Push(e.stack.ReadVar(e.activeScope.bp, index))
}

func (e *emulator) handlePushLocal0() { e.stack.Push(e.stack.ReadVar(e.activeScope.bp, 0)) }
func (e *emulator) handlePushLocal1() { e.stack.Push(e.stack.ReadVar(e.activeScope.bp, 1)) }
func (e *emulator) handlePushLocal2() { e.stack.Push(e.stack.ReadVar(e.activeScope.bp, 2)) }

func (e *emulator) handlePush0()  { e.stack.Push(primitives.Int(0)) }
func (e *emulator) handlePush1()  { e.stack.Push(primitives.Int(1)) }
func (e *emulator) handlePush2()  { e.stack.Push(primitives.Int(2)) }
func (e *emulator) handlePush3()  { e.stack.Push(primitives.Int(3)) }
func (e *emulator) handlePushN1() { e.stack.Push(primitives.Int(-1)) }

func (e *emulator) handlePushInt() {
	value := e.readOperand(instrs.OpPushInt, 0)
	e.stack.Push(primitives.Int(value))
}

func (e *emulator) handlePushStr() {
	index := e.readOperand(instrs.OpPushStr, 0)
	value := primitives.String(e.program.StringPool[index])
	e.stack.Push(value)
}

func (e *emulator) handlePushTrue()  { e.stack.Push(primitives.Bool(true)) }
func (e *emulator) handlePushFalse() { e.stack.Push(primitives.Bool(false)) }
func (e *emulator) handlePushUndef() { e.stack.Push(primitives.Undefined{}) }

func (e *emulator) handlePop() { e.stack.Pop() }

func (e *emulator) handlePopGlobal() {
	index := e.readOperand(instrs.OpPopGlobal, 0)
	e.stack.WriteVar(e.globalScope.bp, index, e.stack.Pop())
}

func (e *emulator) handlePopLocal() {
	index := e.readOperand(instrs.OpPopLocal, 0)
	e.stack.WriteVar(e.activeScope.bp, index, e.stack.Pop())
}

func (e *emulator) handlePopLocal0() { e.stack.WriteVar(e.activeScope.bp, 0, e.stack.Pop()) }
func (e *emulator) handlePopLocal1() { e.stack.WriteVar(e.activeScope.bp, 1, e.stack.Pop()) }
func (e *emulator) handlePopLocal2() { e.stack.WriteVar(e.activeScope.bp, 2, e.stack.Pop()) }

func (e *emulator) handleMakeArray() {
	size := e.readOperand(instrs.OpMakeArray, 0)
	arr := make(primitives.Array, size)
	e.stack.PopMany(arr)
	e.stack.Push(arr)
}

func (e *emulator) handleNewArray() {
	size := e.stack.Pop().(primitives.Int)
	arr := make(primitives.Array, size)
	e.stack.Push(arr)
}

func (e *emulator) handleBinaryAdd() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	err := runtime.InvalidOperandsError{
		Operator: runtime.AddOperator,
		Left:     left,
		Right:    right,
	}

	leftAdder, ok := left.(runtime.Adder)
	if !ok {
		panic(err)
	}

	res, ok := leftAdder.Add(right)
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleBinarySub() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	err := runtime.InvalidOperandsError{
		Operator: runtime.SubOperator,
		Left:     left,
		Right:    right,
	}

	leftSubber, ok := left.(runtime.Subber)
	if !ok {
		panic(err)
	}

	res, ok := leftSubber.Sub(right)
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleBinaryMul() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	err := runtime.InvalidOperandsError{
		Operator: runtime.MulOperator,
		Left:     left,
		Right:    right,
	}

	leftMultiplier, ok := left.(runtime.Multer)
	if !ok {
		panic(err)
	}

	res, ok := leftMultiplier.Mul(right)
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleBinaryDiv() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	err := runtime.InvalidOperandsError{
		Operator: runtime.DivOperator,
		Left:     left,
		Right:    right,
	}

	leftDivider, ok := left.(runtime.Diver)
	if !ok {
		panic(err)
	}

	res, ok := leftDivider.Div(right)
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleBinaryMod() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	err := runtime.InvalidOperandsError{
		Operator: runtime.ModOperator,
		Left:     left,
		Right:    right,
	}

	leftModder, ok := left.(runtime.Modder)
	if !ok {
		panic(err)
	}

	res, ok := leftModder.Mod(right)
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleBinaryEq() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	e.stack.Push(primitives.Bool(left.Equals(right)))
}

func (e *emulator) handleBinaryNeq() {
	right := e.stack.Pop()
	left := e.stack.Pop()

	e.stack.Push(primitives.Bool(!left.Equals(right)))
}

func makeComparisonHandler(op runtime.BinaryOperator) func(*emulator) {
	return func(e *emulator) {
		right := e.stack.Pop()
		left := e.stack.Pop()

		err := runtime.InvalidOperandsError{
			Operator: op,
			Left:     left,
			Right:    right,
		}

		leftComparer, ok := left.(runtime.Comparer)
		if !ok {
			panic(err)
		}

		rel, ok := leftComparer.Compare(right)
		if !ok {
			panic(err)
		}

		var res primitives.Bool
		switch op {
		case runtime.LtOperator:
			res = rel < 0
		case runtime.LtEqOperator:
			res = rel <= 0
		case runtime.GtOperator:
			res = rel > 0
		case runtime.GtEqOperator:
			res = rel >= 0
		default:
			panic("unreachable")
		}

		e.stack.Push(res)
	}
}

func (e *emulator) handleUnaryNeg() {
	operand := e.stack.Pop()

	err := runtime.InvalidOperandError{
		Operator: runtime.NegOperator,
		Operand:  operand,
	}

	negator, ok := operand.(runtime.Negator)
	if !ok {
		panic(err)
	}

	res, ok := negator.Neg()
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleUnaryNot() {
	operand := e.stack.Pop()

	err := runtime.InvalidOperandError{
		Operator: runtime.NotOperator,
		Operand:  operand,
	}

	notter, ok := operand.(runtime.Notter)
	if !ok {
		panic(err)
	}

	res, ok := notter.Not()
	if !ok {
		panic(err)
	}

	e.stack.Push(res)
}

func (e *emulator) handleMakeFunc() {
	templateIndex := e.readOperand(instrs.OpMakeFunc, 0)
	template := e.program.FuncPool[templateIndex]

	captures := make([]*runtime.Obj, len(template.Captures))
	for i, index := range template.Captures {
		// TODO: handle captured globals
		_, _ = i, index
	}

	fn := &bytecode.Func{
		TemplateIndex: templateIndex,
		Captures:      captures,
	}

	e.stack.Push(fn)
}

func (e *emulator) handleCall() {
	numArgs := e.readOperand(instrs.OpCall, 0)

	switch function := e.stack.Pop().(type) {
	case *bytecode.Func:
		template := &e.program.FuncPool[function.TemplateIndex]
		if numArgs > template.NumParams {
			panic(runtime.TooManyArgsError{Expected: template.NumParams, Actual: numArgs})
		}

		bp := e.stack.SP() - numArgs

		e.stack.Grow(template.NumLocals - numArgs)

		newScope := scopePool.Get().(*scope)
		*newScope = scope{
			pc: template.Addr,
			bp: bp,

			function: function,
			parent:   e.activeScope,
		}
		e.activeScope = newScope
	case *builtins.Builtin:
		if function.NumArgs != -1 && numArgs != function.NumArgs {
			panic(runtime.TooManyArgsError{Expected: function.NumArgs, Actual: numArgs})
		}

		args := make([]runtime.Obj, numArgs)
		e.stack.PopMany(args)

		res, err := function.Func(args)
		if err != nil {
			panic(err)
		}

		e.stack.Push(res)
	default:
		panic(runtime.NotCallableError{Func: function})
	}
}

func (e *emulator) handleReturn() {
	restoredScope := e.activeScope
	e.activeScope = e.activeScope.parent

	e.stack.PatchReturn(restoredScope.bp)

	*restoredScope = scope{}
	scopePool.Put(restoredScope)
}

func (_ *emulator) handleHalt() {
	panic(struct{}{})
}

func (e *emulator) handlePushCollElem() {
	key := e.stack.Pop()
	coll := e.stack.Pop()

	err := runtime.NotIndexableError{
		Coll: coll,
		Key:  key,
	}

	indexer, ok := coll.(runtime.Coll)
	if !ok {
		panic(err)
	}

	elem, ok := indexer.Get(key)
	if !ok {
		panic(err)
	}

	e.stack.Push(elem)
}

func (e *emulator) handlePopCollElem() {
	value := e.stack.Pop()
	key := e.stack.Pop()
	coll := e.stack.Pop()

	err := runtime.NotIndexableError{
		Coll: coll,
		Key:  key,
	}

	indexer, ok := coll.(runtime.Coll)
	if !ok {
		panic(err)
	}

	ok = indexer.Set(key, value)
	if !ok {
		panic(err)
	}
}
