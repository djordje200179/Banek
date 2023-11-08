package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode"
	"banek/bytecode/instructions"
	"banek/compiler/scopes"
	"banek/exec/errors"
	"banek/exec/objects"
)

func (compiler *compiler) compileExpr(expr ast.Expression) error {
	scope := compiler.topScope()

	switch expr := expr.(type) {
	case expressions.ConstLiteral:
		scope.EmitInstr(instructions.OpPushConst, compiler.addConst(expr.Value))
		return nil
	case expressions.BinaryOp:
		return compiler.compileBinaryOp(expr)
	case expressions.UnaryOp:
		return compiler.compileUnaryOp(expr)
	case expressions.If:
		return compiler.compileIfExpr(expr)
	case expressions.ArrayLiteral:
		for _, elem := range expr {
			err := compiler.compileExpr(elem)
			if err != nil {
				return err
			}
		}

		scope.EmitInstr(instructions.OpNewArray, len(expr))

		return nil
	case expressions.CollIndex:
		err := compiler.compileExpr(expr.Coll)
		if err != nil {
			return err
		}

		err = compiler.compileExpr(expr.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instructions.OpPushCollElem)

		return nil
	case expressions.Assignment:
		return compiler.compileAssigment(expr)
	case expressions.FuncCall:
		for _, arg := range expr.Args {
			err := compiler.compileExpr(arg)
			if err != nil {
				return err
			}
		}

		err := compiler.compileExpr(expr.Func)
		if err != nil {
			return err
		}

		scope.EmitInstr(instructions.OpCall, len(expr.Args))

		return nil
	case expressions.FuncLiteral:
		return compiler.compileFuncLiteral(expr)
	case expressions.Identifier:
		return compiler.compileIdentifier(expr)
	default:
		return ast.ErrUnknownExpr{Expr: expr}
	}
}

func (compiler *compiler) compileIfExpr(expression expressions.If) error {
	err := compiler.compileExpr(expression.Cond)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instructions.OpBranchIfFalse, 0)

	err = compiler.compileExpr(expression.Consequence)
	if err != nil {
		return err
	}

	branchSize := instructions.OpBranch.Info().Size()
	branchIfFalseSize := instructions.OpBranchIfFalse.Info().Size()

	secondPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instructions.OpBranch, 0)
	elseAddr := secondPatchAddr + branchSize

	err = compiler.compileExpr(expression.Alternative)
	if err != nil {
		return err
	}

	outAddr := scope.CurrAddr()

	scope.PatchInstrOperand(secondPatchAddr, 0, outAddr-secondPatchAddr-branchSize)
	scope.PatchInstrOperand(firstPatchAddr, 0, elseAddr-firstPatchAddr-branchIfFalseSize)

	return nil
}

func (compiler *compiler) compileAssigment(expr expressions.Assignment) error {
	err := compiler.compileExpr(expr.Value)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	scope.EmitInstr(instructions.OpPushDup)

	switch variable := expr.Var.(type) {
	case expressions.Identifier:
		varName := variable.String()

		var varScope scopes.Scope
		var varIndex, varScopeIndex int
		for i := len(compiler.scopes) - 1; i >= 0; i-- {
			variable, index := compiler.scopes[i].GetVar(varName)
			if index == -1 {
				continue
			}

			if !variable.Mutable {
				return errors.ErrIdentifierNotMutable{Identifier: varName}
			}

			varScope = compiler.scopes[i]
			varScopeIndex = i
			varIndex = index

			break
		}

		if varScope == nil {
			return errors.ErrIdentifierNotDefined{Identifier: varName}
		}

		if varScope.IsGlobal() {
			scope.EmitInstr(instructions.OpPopGlobal, varIndex)
			return nil
		} else if varScope == scope {
			scope.EmitInstr(instructions.OpPopLocal, varIndex)
			return nil
		}

		varScope.MarkCaptured()

		capturedVarLevel := len(compiler.scopes) - 2 - varScopeIndex

		funcScope, ok := varScope.(*scopes.Function)
		if !ok {
			block := varScope.(*scopes.Block)
			for {
				nextScope := block.Parent

				funcScope, ok = nextScope.(*scopes.Function)
				if ok {
					break
				}

				block = nextScope.(*scopes.Block)
			}
		}

		capturedVarIndex := funcScope.AddCapturedVar(capturedVarLevel, varIndex)
		scope.EmitInstr(instructions.OpPopCaptured, capturedVarIndex)

		return nil
	case expressions.CollIndex:
		err := compiler.compileExpr(variable.Coll)
		if err != nil {
			return err
		}

		err = compiler.compileExpr(variable.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instructions.OpPopCollElem)

		return nil
	default:
		return ast.ErrInvalidAssignment{Variable: expr.Var}
	}
}

func (compiler *compiler) compileFuncLiteral(expr expressions.FuncLiteral) error {
	funcScope := new(scopes.Function)

	paramNames := make([]string, len(expr.Params))
	for i, param := range expr.Params {
		paramNames[i] = param.String()
	}

	err := funcScope.AddParams(paramNames)
	if err != nil {
		return err
	}

	compiler.pushScope(funcScope)
	err = compiler.compileExpr(expr.Body)
	if err != nil {
		return err
	}
	funcScope.EmitInstr(instructions.OpReturn)
	compiler.popScope()

	funcTemplate := funcScope.MakeFunction()

	funcIndex := compiler.addFunc(funcTemplate)

	scope := compiler.topScope()

	if funcTemplate.IsClosure() {
		scope.EmitInstr(instructions.OpNewFunc, funcIndex)
	} else {
		funcObject := &bytecode.Func{
			TemplateIndex: funcIndex,
		}

		scope.EmitInstr(instructions.OpPushConst, compiler.addConst(funcObject))
	}

	return nil
}

func (compiler *compiler) compileIdentifier(expr expressions.Identifier) error {
	varName := expr.String()

	scope := compiler.topScope()

	if index := objects.BuiltinFindIndex(varName); index != -1 {
		scope.EmitInstr(instructions.OpPushBuiltin, index)
		return nil
	}

	var varScope scopes.Scope
	var varIndex, varScopeIndex int
	for i := len(compiler.scopes) - 1; i >= 0; i-- {
		_, index := compiler.scopes[i].GetVar(varName)
		if index == -1 {
			continue
		}

		varScope = compiler.scopes[i]
		varScopeIndex = i
		varIndex = index

		break
	}

	if varScope == nil {
		return errors.ErrIdentifierNotDefined{Identifier: varName}
	}

	if varScope.IsGlobal() {
		scope.EmitInstr(instructions.OpPushGlobal, varIndex)
		return nil
	} else if varScope == scope {
		scope.EmitInstr(instructions.OpPushLocal, varIndex)
		return nil
	}

	varScope.MarkCaptured()

	capturedVarLevel := len(compiler.scopes) - 2 - varScopeIndex

	funcScope, ok := scope.(*scopes.Function)
	if !ok {
		block := scope.(*scopes.Block)
		for {
			nextScope := block.Parent

			funcScope, ok = nextScope.(*scopes.Function)
			if ok {
				break
			}

			block = nextScope.(*scopes.Block)
		}
	}

	capturedVarIndex := funcScope.AddCapturedVar(capturedVarLevel, varIndex)
	scope.EmitInstr(instructions.OpPushCaptured, capturedVarIndex)

	return nil
}

func (compiler *compiler) compileBinaryOp(expression expressions.BinaryOp) error {
	err := compiler.compileExpr(expression.Left)
	if err != nil {
		return err
	}

	err = compiler.compileExpr(expression.Right)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instructions.OpBinaryOp, int(expression.Operator))

	return nil
}

func (compiler *compiler) compileUnaryOp(expression expressions.UnaryOp) error {
	err := compiler.compileExpr(expression.Operand)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instructions.OpUnaryOp, int(expression.Operator))

	return nil
}
