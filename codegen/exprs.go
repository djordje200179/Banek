package codegen

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/symtable/symbols"
	"banek/tokens"
	"slices"
)

func (g *generator) compileExpr(expr ast.Expr) {
	switch expr := expr.(type) {
	case exprs.ArrayLiteral:
		g.compileArrayLiteral(expr)
	case exprs.BinaryOp:
		g.compileBinaryOp(expr)
	case exprs.BoolLiteral:
		g.compileBoolLiteral(expr)
	case exprs.CollIndex:
		g.compileCollIndex(expr, true)
	case exprs.FuncCall:
		g.compileFuncCall(expr)
	case exprs.FuncLiteral:
		g.compileFuncLiteral(expr)
	case exprs.Ident:
		g.compileIdent(expr)
	case exprs.If:
		g.compileIfExpr(expr)
	case exprs.IntLiteral:
		g.compileIntLiteral(expr)
	case exprs.StringLiteral:
		g.compileStringLiteral(expr)
	case exprs.UnaryOp:
		g.compileUnaryOp(expr)
	case exprs.UndefinedLiteral:
		g.compileUndefinedLiteral(expr)
	default:
		panic("unreachable")
	}
}

func (g *generator) compileArrayLiteral(expr exprs.ArrayLiteral) {
	for _, elem := range expr {
		g.compileExpr(elem)
	}

	g.emitInstr(instrs.OpMakeArray, len(expr))
}

func (g *generator) compileBinaryOp(expr exprs.BinaryOp) {
	g.compileExpr(expr.Left)
	g.compileExpr(expr.Right)

	switch expr.Operator {
	case tokens.Plus:
		g.emitInstr(instrs.OpAdd)
	case tokens.Minus:
		g.emitInstr(instrs.OpSub)
	case tokens.Asterisk:
		g.emitInstr(instrs.OpMul)
	case tokens.Slash:
		g.emitInstr(instrs.OpDiv)
	case tokens.Percent:
		g.emitInstr(instrs.OpMod)
	case tokens.Equals:
		g.emitInstr(instrs.OpCompareEq)
	case tokens.NotEquals:
		g.emitInstr(instrs.OpCompareNeq)
	case tokens.Greater:
		g.emitInstr(instrs.OpCompareGt)
	case tokens.GreaterEquals:
		g.emitInstr(instrs.OpCompareGtEq)
	case tokens.Less:
		g.emitInstr(instrs.OpCompareLt)
	case tokens.LessEquals:
		g.emitInstr(instrs.OpCompareLtEq)
	default:
		panic("unreachable")
	}

}

func (g *generator) compileCollIndex(expr exprs.CollIndex, load bool) {
	g.compileExpr(expr.Coll)
	g.compileExpr(expr.Key)

	if load {
		g.emitInstr(instrs.OpCollGet)
	}
}

func (g *generator) compileFuncCall(expr exprs.FuncCall) {
	for _, arg := range expr.Args {
		g.compileExpr(arg)
	}

	g.compileExpr(expr.Func)
	g.emitInstr(instrs.OpCall, len(expr.Args))
}

func (g *generator) compileIdent(expr exprs.Ident) {
	switch sym := expr.Symbol.(type) {
	case symbols.Builtin:
		g.emitInstr(instrs.OpBuiltin, int(sym))
	case symbols.Var:
		switch {
		case sym.Level == 0:
			g.emitInstr(instrs.OpLoadGlobal, sym.Index)
		case sym.Level == g.active.level:
			switch sym.Index {
			case 0:
				g.emitInstr(instrs.OpLoadLocal0)
			case 1:
				g.emitInstr(instrs.OpLoadLocal1)
			case 2:
				g.emitInstr(instrs.OpLoadLocal2)
			default:
				g.emitInstr(instrs.OpLoadLocal, sym.Index)
			}
		default:
			g.emitInstr(instrs.OpLoadCaptured, g.active.level-sym.Level, sym.Index)
		}
	default:
		panic("unreachable")
	}
}

func (g *generator) compileIfExpr(expr exprs.If) {
	g.compileExpr(expr.Cond)
	jumpPC := g.currAddr()
	g.emitInstr(instrs.OpBranchFalse, 0)

	g.compileExpr(expr.Cons)

	altJumpPC := g.currAddr()
	g.emitInstr(instrs.OpJump, 0)
	g.patchJumpOperand(jumpPC, 0)
	g.compileExpr(expr.Alt)
	g.patchJumpOperand(altJumpPC, 0)
}

func (g *generator) compileUnaryOp(expr exprs.UnaryOp) {
	g.compileExpr(expr.Operand)

	switch expr.Operator {
	case tokens.Bang:
		g.emitInstr(instrs.OpNot)
	case tokens.Minus:
		g.emitInstr(instrs.OpNeg)
	default:
		panic("unreachable")
	}
}

func (g *generator) compileBoolLiteral(expr exprs.BoolLiteral) {
	if expr {
		g.emitInstr(instrs.OpConstTrue)
	} else {
		g.emitInstr(instrs.OpConstFalse)
	}
}

func (g *generator) compileIntLiteral(expr exprs.IntLiteral) {
	switch int(expr) {
	case 0:
		g.emitInstr(instrs.OpConst0)
	case 1:
		g.emitInstr(instrs.OpConst1)
	case 2:
		g.emitInstr(instrs.OpConst2)
	case 3:
		g.emitInstr(instrs.OpConst3)
	case -1:
		g.emitInstr(instrs.OpConstN1)
	default:
		g.emitInstr(instrs.OpConstInt, int(expr))
	}
}

func (g *generator) compileStringLiteral(expr exprs.StringLiteral) {
	i := slices.Index(g.stringPool, string(expr))
	if i == -1 {
		i = len(g.stringPool)
		g.stringPool = append(g.stringPool, string(expr))
	}

	g.emitInstr(instrs.OpConstStr, i)
}

func (g *generator) compileUndefinedLiteral(_ exprs.UndefinedLiteral) {
	g.emitInstr(instrs.OpConstUndef)
}

func (g *generator) compileFuncLiteral(expr exprs.FuncLiteral) {
	funcIndex := len(g.funcPool)
	g.active = &container{
		level:    g.active.level + 1,
		index:    funcIndex,
		previous: g.active,
		vars:     len(expr.Params),
	}

	f := bytecode.Func{
		NumParams: len(expr.Params),
	}

	g.compileExpr(expr.Body)
	g.emitInstr(instrs.OpReturn)

	f.NumLocals = g.active.vars
	g.funcPool = append(g.funcPool, f)
	g.funcCodes[funcIndex] = g.active.code

	g.active = g.active.previous

	g.emitInstr(instrs.OpMakeFunc, funcIndex)
}

func (g *generator) compileStore(expr ast.Expr) {
	switch expr := expr.(type) {
	case exprs.Ident:
		v := expr.Symbol.(symbols.Var)

		switch {
		case v.Level == 0:
			g.emitInstr(instrs.OpStoreGlobal, v.Index)
		case v.Level == g.active.level:
			switch v.Index {
			case 0:
				g.emitInstr(instrs.OpStoreLocal0)
			case 1:
				g.emitInstr(instrs.OpStoreLocal1)
			case 2:
				g.emitInstr(instrs.OpStoreLocal2)
			default:
				g.emitInstr(instrs.OpStoreLocal, v.Index)
			}
		default:
			g.emitInstr(instrs.OpStoreCaptured, g.active.level-v.Level, v.Index)
		}
	case exprs.CollIndex:
		g.emitInstr(instrs.OpCollSet)
	default:
		panic("unreachable")
	}
}

func (g *generator) compilePreStore(expr ast.Expr) {
	switch expr := expr.(type) {
	case exprs.Ident:
	case exprs.CollIndex:
		g.compileCollIndex(expr, false)
	default:
		panic("unreachable")
	}
}
