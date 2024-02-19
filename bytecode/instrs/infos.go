package instrs

var instrInfos = [...]InstrInfo{
	OpHalt:   {"HALT", nil},
	OpCall:   {"CALL", []OperandInfo{{1, OperandLiteral}}},
	OpReturn: {"RET", nil},

	OpJump:        {"JMP", []OperandInfo{{2, OperandOffset}}},
	OpBranchFalse: {"BR.F", []OperandInfo{{2, OperandOffset}}},
	OpBranchTrue:  {"BR.T", []OperandInfo{{2, OperandOffset}}},

	OpPushBuiltin:  {"PUSH.B", []OperandInfo{{1, OperandBuiltin}}},
	OpPushGlobal:   {"PUSH.G", []OperandInfo{{2, OperandLiteral}}},
	OpPushCaptured: {"PUSH.C", []OperandInfo{{1, OperandLiteral}, {1, OperandLiteral}}},
	OpPushCollElem: {"PUSH.CE", nil},
	OpPushLocal:    {"PUSH.L", []OperandInfo{{1, OperandLiteral}}},
	OpPushLocal0:   {"PUSH.L0", nil},
	OpPushLocal1:   {"PUSH.L1", nil},
	OpPushLocal2:   {"PUSH.L2", nil},

	OpPush0:     {"PUSH.0", nil},
	OpPush1:     {"PUSH.1", nil},
	OpPush2:     {"PUSH.2", nil},
	OpPush3:     {"PUSH.3", nil},
	OpPushN1:    {"PUSH.N1", nil},
	OpPushInt:   {"PUSH.I", []OperandInfo{{8, OperandLiteral}}},
	OpPushStr:   {"PUSH.S", []OperandInfo{{2, OperandLiteral}}},
	OpPushTrue:  {"PUSH.TRUE", nil},
	OpPushFalse: {"PUSH.FALSE", nil},
	OpPushUndef: {"PUSH.UNDEF", nil},

	OpPop:         {"POP", nil},
	OpPopGlobal:   {"POP.G", []OperandInfo{{2, OperandLiteral}}},
	OpPopCaptured: {"POP.C", []OperandInfo{{1, OperandLiteral}, {1, OperandLiteral}}},
	OpPopCollElem: {"POP.CE", nil},
	OpPopLocal:    {"POP.L", []OperandInfo{{1, OperandLiteral}}},
	OpPopLocal0:   {"POP.L0", nil},
	OpPopLocal1:   {"POP.L1", nil},
	OpPopLocal2:   {"POP.L2", nil},

	OpDup:  {"DUP.1", nil},
	OpDup2: {"DUP.2", nil},
	OpSwap: {"SWAP", nil},

	OpBinaryAdd: {"ADD", nil},
	OpBinarySub: {"SUB", nil},
	OpBinaryMul: {"MUL", nil},
	OpBinaryDiv: {"DIV", nil},
	OpBinaryMod: {"MOD", nil},
	OpBinaryEq:  {"CMP.EQ", nil},
	OpBinaryNe:  {"CMP.NE", nil},
	OpBinaryGt:  {"CMP.GT", nil},
	OpBinaryGe:  {"CMP.GE", nil},
	OpBinaryLt:  {"CMP.LT", nil},
	OpBinaryLe:  {"CMP.LE", nil},

	OpUnaryNeg: {"NEG", nil},
	OpUnaryNot: {"NOT", nil},

	OpMakeArray: {"MAKE.ARRAY", []OperandInfo{{1, OperandLiteral}}},
	OpNewArray:  {"NEW.ARRAY", nil},
	OpMakeFunc:  {"MAKE.FUNC", []OperandInfo{{2, OperandFunc}}},
}
