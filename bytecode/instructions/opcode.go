package instructions

type Opcode byte

const (
	OpInvalid Opcode = iota

	OpPushDup
	OpPushConst
	OpPushLocal
	OpPushGlobal
	OpPushCaptured
	OpPushBuiltin
	OpPushCollElem

	OpPop
	OpPopLocal
	OpPopGlobal
	OpPopCaptured
	OpPopCollElem

	OpBinaryOp
	OpUnaryOp

	OpBranch
	OpBranchIfFalse

	OpCall
	OpReturn

	OpNewArray
	OpNewFunc
)

func (opcode Opcode) String() string {
	return opcode.Info().Name
}

func (opcode Opcode) Info() InstrInfo {
	if opcode < 0 || opcode >= Opcode(len(InstrInfos)) {
		return InstrInfos[OpInvalid]
	}

	return InstrInfos[opcode]
}

type InstrInfo struct {
	Name string

	Operands []OperandInfo
}

func (instrInfo InstrInfo) Size() int {
	size := 1
	for _, operand := range instrInfo.Operands {
		size += operand.Width
	}

	return size
}

func (instrInfo InstrInfo) OperandOffset(index int) int {
	offset := 1
	for i := 0; i < index; i++ {
		offset += instrInfo.Operands[i].Width
	}

	return offset
}

var InstrInfos = [...]InstrInfo{
	OpInvalid: {"INVALID", []OperandInfo{}},

	OpPushDup:      {"PUSH.D", []OperandInfo{}},
	OpPushConst:    {"PUSH.C", []OperandInfo{{2, OperandConstant}}},
	OpPushLocal:    {"PUSH.L", []OperandInfo{{1, OperandLiteral}}},
	OpPushGlobal:   {"PUSH.G", []OperandInfo{{1, OperandLiteral}}},
	OpPushCaptured: {"PUSH.O", []OperandInfo{{1, OperandLiteral}}},
	OpPushBuiltin:  {"PUSH.B", []OperandInfo{{1, OperandLiteral}}},
	OpPushCollElem: {"PUSH.CE", []OperandInfo{}},

	OpPop:         {"POP", []OperandInfo{}},
	OpPopLocal:    {"POP.L", []OperandInfo{{1, OperandLiteral}}},
	OpPopGlobal:   {"POP.G", []OperandInfo{{1, OperandLiteral}}},
	OpPopCaptured: {"POP.O", []OperandInfo{{1, OperandLiteral}}},
	OpPopCollElem: {"POP.CE", []OperandInfo{}},

	OpBinaryOp: {"OP.I", []OperandInfo{{1, OperandInfixOp}}},
	OpUnaryOp:  {"OP.P", []OperandInfo{{1, OperandPrefixOp}}},

	OpBranch:        {"BR", []OperandInfo{{2, OperandLiteral}}},
	OpBranchIfFalse: {"BR.F", []OperandInfo{{2, OperandLiteral}}},

	OpCall:   {"CALL", []OperandInfo{{1, OperandLiteral}}},
	OpReturn: {"RET", []OperandInfo{}},

	OpNewArray: {"NEW.A", []OperandInfo{{2, OperandLiteral}}},
	OpNewFunc:  {"NEW.F", []OperandInfo{{2, OperandFunc}}},
}