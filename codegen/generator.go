package codegen

import (
	"banek/ast"
	"banek/bytecode"
	"banek/bytecode/instrs"
	"slices"
)

func Generate(stmtChan <-chan ast.Stmt) bytecode.Executable {
	g := &generator{
		funcPool:  make([]bytecode.FuncTemplate, 1),
		funcCodes: make(map[int]instrs.Code),
	}
	g.active = &g.global

	for stmt := range stmtChan {
		g.compileStmt(stmt)
	}

	g.emitInstr(instrs.OpHalt)

	g.funcPool[0] = bytecode.FuncTemplate{
		Name:       "<entry>",
		NumLocals:  g.global.vars,
		IsCaptured: true,
	}

	return g.makeExecutable()
}

type container struct {
	level, index int
	vars         int

	code instrs.Code

	previous *container
}

type generator struct {
	global container
	active *container

	funcCodes map[int]instrs.Code

	stringPool []string
	funcPool   []bytecode.FuncTemplate
}

func (g *generator) makeExecutable() bytecode.Executable {
	totalCodeSize := len(g.global.code)
	for _, code := range g.funcCodes {
		totalCodeSize += len(code)
	}

	code := make(instrs.Code, 0, totalCodeSize)
	code = append(code, g.global.code...)
	for funcIndex, funcCode := range g.funcCodes {
		g.funcPool[funcIndex].Addr = len(code)
		code = append(code, funcCode...)
	}

	return bytecode.Executable{
		Code:       code,
		StringPool: g.stringPool,
		FuncPool:   g.funcPool,
	}
}

func (g *generator) emitInstr(opcode instrs.Opcode, operands ...int) {
	instr := instrs.MakeInstr(opcode, operands...)
	g.active.code = slices.Concat(g.active.code, instr)
}

func (g *generator) patchJumpOperand(addr int, operandIndex int) {
	op := instrs.Opcode(g.active.code[addr])
	opInfo := op.Info()

	instCode := g.active.code[addr : addr+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	offset := g.currAddr() - addr - opInfo.Size()
	copy(instCode[operandOffset:], instrs.MakeOperandValue(offset, operandWidth))
}

func (g *generator) currAddr() int {
	return len(g.active.code)
}
