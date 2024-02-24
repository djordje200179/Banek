package instrs

var instrInfos = [...]InstrInfo{
	OpHalt:   {"HALT", nil},
	OpCall:   {"CALL", []OperandInfo{{1, OperandLiteral}}},
	OpReturn: {"RET", nil},

	OpJump:        {"JMP", []OperandInfo{{2, OperandOffset}}},
	OpBranchFalse: {"BR.F", []OperandInfo{{2, OperandOffset}}},
	OpBranchTrue:  {"BR.T", []OperandInfo{{2, OperandOffset}}},

	OpLoadGlobal:   {"LD.G", []OperandInfo{{2, OperandLiteral}}},
	OpLoadCaptured: {"LD.C", []OperandInfo{{1, OperandLiteral}, {1, OperandLiteral}}},
	OpLoadLocal:    {"LD.L", []OperandInfo{{1, OperandLiteral}}},
	OpLoadLocal0:   {"LD.L0", nil},
	OpLoadLocal1:   {"LD.L1", nil},
	OpLoadLocal2:   {"LD.L2", nil},

	OpConst0:     {"CONST.0", nil},
	OpConst1:     {"CONST.1", nil},
	OpConst2:     {"CONST.2", nil},
	OpConst3:     {"CONST.3", nil},
	OpConstN1:    {"CONST.N1", nil},
	OpConstInt:   {"CONST.I", []OperandInfo{{8, OperandLiteral}}},
	OpConstStr:   {"CONST.S", []OperandInfo{{2, OperandLiteral}}},
	OpConstTrue:  {"CONST.T", nil},
	OpConstFalse: {"CONST.F", nil},
	OpConstUndef: {"CONST.U", nil},
	OpBuiltin:    {"BUILTIN", []OperandInfo{{1, OperandLiteral}}},

	OpMakeArray: {"MAKE.ARRAY", []OperandInfo{{1, OperandLiteral}}},
	OpNewArray:  {"NEW.ARRAY", nil},
	OpMakeFunc:  {"MAKE.FUNC", []OperandInfo{{2, OperandFunc}}},

	OpPop:           {"POP", nil},
	OpStoreGlobal:   {"ST.G", []OperandInfo{{2, OperandLiteral}}},
	OpStoreCaptured: {"ST.C", []OperandInfo{{1, OperandLiteral}, {1, OperandLiteral}}},
	OpStoreLocal:    {"ST.L", []OperandInfo{{1, OperandLiteral}}},
	OpStoreLocal0:   {"ST.L0", nil},
	OpStoreLocal1:   {"ST.L1", nil},
	OpStoreLocal2:   {"ST.L2", nil},

	OpCollGet: {"COLL.GET", nil},
	OpCollSet: {"COLL.SET", nil},

	OpDup:  {"DUP", nil},
	OpDup2: {"DUP2", nil},
	OpSwap: {"SWAP", nil},

	OpCompareEq:   {"CMP.EQ", nil},
	OpCompareNeq:  {"CMP.NE", nil},
	OpCompareGt:   {"CMP.GT", nil},
	OpCompareGtEq: {"CMP.GE", nil},
	OpCompareLt:   {"CMP.LT", nil},
	OpCompareLtEq: {"CMP.LE", nil},

	OpAdd: {"ADD", nil},
	OpSub: {"SUB", nil},
	OpMul: {"MUL", nil},
	OpDiv: {"DIV", nil},
	OpMod: {"MOD", nil},

	OpNeg: {"NEG", nil},
	OpNot: {"NOT", nil},
}
