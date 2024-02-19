package emulator

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime"
	"banek/runtime/builtins"
	"banek/runtime/primitives"
)

func (e *emulator) handleDup()  { e.operandStack.Dup() }
func (e *emulator) handleDup2() { e.operandStack.Dup2() }
func (e *emulator) handleDup3() { e.operandStack.Dup3() }
func (e *emulator) handleSwap() { e.operandStack.Swap() }

func (e *emulator) handleJump() {
	offset := e.scopeStack.ReadOperand(instrs.OpJump, 0)
	e.scopeStack.MovePC(offset)
}

func (e *emulator) handleBranchFalse() {
	offset := e.scopeStack.ReadOperand(instrs.OpBranchFalse, 0)

	if !e.operandStack.Pop().Truthy() {
		e.scopeStack.MovePC(offset)
	}
}

func (e *emulator) handleBranchTrue() {
	offset := e.scopeStack.ReadOperand(instrs.OpBranchTrue, 0)

	if e.operandStack.Pop().Truthy() {
		e.scopeStack.MovePC(offset)
	}
}

func (e *emulator) handlePushBuiltin() {
	builtin := e.scopeStack.ReadOperand(instrs.OpPushBuiltin, 0)
	e.operandStack.Push(&builtins.Funcs[builtin])
}

func (e *emulator) handlePushGlobal() {
	index := e.scopeStack.ReadOperand(instrs.OpPushGlobal, 0)
	e.operandStack.Push(e.scopeStack.GetGlobal(index))
}

func (e *emulator) handlePushLocal() {
	index := e.scopeStack.ReadOperand(instrs.OpPushLocal, 0)
	e.operandStack.Push(e.scopeStack.GetLocal(index))
}

func (e *emulator) handlePushLocal0() { e.operandStack.Push(e.scopeStack.GetLocal(0)) }
func (e *emulator) handlePushLocal1() { e.operandStack.Push(e.scopeStack.GetLocal(1)) }
func (e *emulator) handlePushLocal2() { e.operandStack.Push(e.scopeStack.GetLocal(2)) }

func (e *emulator) handlePush0()  { e.operandStack.Push(primitives.Int(0)) }
func (e *emulator) handlePush1()  { e.operandStack.Push(primitives.Int(1)) }
func (e *emulator) handlePush2()  { e.operandStack.Push(primitives.Int(2)) }
func (e *emulator) handlePush3()  { e.operandStack.Push(primitives.Int(3)) }
func (e *emulator) handlePushN1() { e.operandStack.Push(primitives.Int(-1)) }

func (e *emulator) handlePushInt() {
	value := e.scopeStack.ReadOperand(instrs.OpPushInt, 0)
	e.operandStack.Push(primitives.Int(value))
}

func (e *emulator) handlePushStr() {
	index := e.scopeStack.ReadOperand(instrs.OpPushStr, 0)
	value := primitives.String(e.program.StringPool[index])
	e.operandStack.Push(value)
}

func (e *emulator) handlePushTrue()  { e.operandStack.Push(primitives.Bool(true)) }
func (e *emulator) handlePushFalse() { e.operandStack.Push(primitives.Bool(false)) }
func (e *emulator) handlePushUndef() { e.operandStack.Push(primitives.Undefined{}) }

func (e *emulator) handlePop() { e.operandStack.Pop() }

func (e *emulator) handlePopGlobal() {
	index := e.scopeStack.ReadOperand(instrs.OpPopGlobal, 0)
	e.scopeStack.SetGlobal(index, e.operandStack.Pop())
}

func (e *emulator) handlePopLocal() {
	index := e.scopeStack.ReadOperand(instrs.OpPopLocal, 0)
	e.scopeStack.SetLocal(index, e.operandStack.Pop())
}

func (e *emulator) handlePopLocal0() { e.scopeStack.SetLocal(0, e.operandStack.Pop()) }
func (e *emulator) handlePopLocal1() { e.scopeStack.SetLocal(1, e.operandStack.Pop()) }
func (e *emulator) handlePopLocal2() { e.scopeStack.SetLocal(2, e.operandStack.Pop()) }

func (e *emulator) handleMakeArray() {
	size := e.scopeStack.ReadOperand(instrs.OpMakeArray, 0)
	arr := make(primitives.Array, size)
	e.operandStack.PopMany(arr)
	e.operandStack.Push(arr)
}

func (e *emulator) handleNewArray() {
	size := e.operandStack.Pop().(primitives.Int)
	arr := make(primitives.Array, size)
	e.operandStack.Push(arr)
}

func (e *emulator) handleBinaryAdd() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleBinarySub() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleBinaryMul() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleBinaryDiv() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleBinaryMod() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleBinaryEq() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

	e.operandStack.Push(primitives.Bool(left.Equals(right)))
}

func (e *emulator) handleBinaryNeq() {
	right := e.operandStack.Pop()
	left := e.operandStack.Pop()

	e.operandStack.Push(primitives.Bool(!left.Equals(right)))
}

func makeComparisonHandler(op runtime.BinaryOperator) func(*emulator) {
	return func(e *emulator) {
		right := e.operandStack.Pop()
		left := e.operandStack.Pop()

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

		e.operandStack.Push(res)
	}
}

func (e *emulator) handleUnaryNeg() {
	operand := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleUnaryNot() {
	operand := e.operandStack.Pop()

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

	e.operandStack.Push(res)
}

func (e *emulator) handleMakeFunc() {
	templateIndex := e.scopeStack.ReadOperand(instrs.OpMakeFunc, 0)
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

	e.operandStack.Push(fn)
}

func (e *emulator) handleCall() {
	numArgs := e.scopeStack.ReadOperand(instrs.OpCall, 0)

	switch function := e.operandStack.Pop().(type) {
	case *bytecode.Func:
		template := &e.program.FuncPool[function.TemplateIndex]
		if numArgs > template.NumParams {
			panic(runtime.TooManyArgsError{Expected: template.NumParams, Actual: numArgs})
		}

		locals := e.scopeStack.NewScope(function, template)
		e.operandStack.PopMany(locals[:numArgs])
	case *builtins.Builtin:
		if function.NumArgs != -1 && numArgs != function.NumArgs {
			panic(runtime.TooManyArgsError{Expected: function.NumArgs, Actual: numArgs})
		}

		args := make([]runtime.Obj, numArgs)
		e.operandStack.PopMany(args)

		res, err := function.Func(args)
		if err != nil {
			panic(err)
		}

		e.operandStack.Push(res)
	default:
		panic(runtime.NotCallableError{Func: function})
	}
}

func (e *emulator) handleReturn() {
	e.scopeStack.RestoreScope()
}

func (_ *emulator) handleHalt() {
	panic(struct{}{})
}
