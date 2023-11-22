package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/objs"
	"unsafe"
)

type Stack struct {
	global Scope
	active Scope

	last *Scope
}

func NewStack(program bytecode.Executable) *Stack {
	globals := make([]objs.Obj, program.NumGlobals)

	stack := &Stack{
		global: Scope{
			vars: globals,
			code: program.Code,
		},
	}
	stack.active = stack.global

	return stack
}

func (stack *Stack) GetGlobal(index int) objs.Obj        { return stack.global.vars[index] }
func (stack *Stack) SetGlobal(index int, value objs.Obj) { stack.global.vars[index] = value }

func (stack *Stack) GetLocal(index int) objs.Obj        { return stack.active.vars[index] }
func (stack *Stack) SetLocal(index int, value objs.Obj) { stack.active.vars[index] = value }

func (stack *Stack) GetCaptured(index int) objs.Obj { return *stack.active.function.Captures[index] }
func (stack *Stack) SetCaptured(index int, value objs.Obj) {
	*stack.active.function.Captures[index] = value
}

func (stack *Stack) GetCapture(captureInfo bytecode.Capture) *objs.Obj {
	varScope := stack.last
	for j := 1; j < captureInfo.Level; j++ {
		varScope = varScope.parent
	}

	return &varScope.vars[captureInfo.Index]
}

func (stack *Stack) NewScope(function *bytecode.Func, template *bytecode.FuncTemplate) []objs.Obj {
	funcScope := scopePool.Get().(*Scope)
	*funcScope = stack.active
	stack.last = funcScope

	locals := newScopeVars(template.NumLocals)

	stack.active = Scope{
		code:     template.Code,
		vars:     locals,
		function: function,
		template: template,
		parent:   stack.last,
	}

	return locals
}

func (stack *Stack) RestoreScope() {
	if !stack.active.template.IsCaptured {
		freeScopeVars(stack.active.vars)
	}

	restoredScope := stack.last
	stack.last = restoredScope.parent
	stack.active = *restoredScope

	*restoredScope = Scope{}
	scopePool.Put(restoredScope)
}

func (stack *Stack) ReadOpcode() instrs.Opcode {
	codePtr := unsafe.Pointer(unsafe.SliceData(stack.active.code))
	offset := uintptr(stack.active.pc) * unsafe.Sizeof(instrs.Opcode(0))
	addr := unsafe.Add(codePtr, offset)

	opcode := *(*instrs.Opcode)(addr)
	stack.active.pc++

	return opcode
}

func (stack *Stack) ReadOperand(width int) int {
	operandSlice := stack.active.code[stack.active.pc : stack.active.pc+width]
	operandValue := instrs.ReadOperandValue(operandSlice)
	stack.active.pc += width

	return operandValue
}

func (stack *Stack) MovePC(offset int) {
	stack.active.pc += offset
}
