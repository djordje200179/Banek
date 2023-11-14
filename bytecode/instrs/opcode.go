package instrs

type Opcode byte

const (
	OpInvalid Opcode = iota

	OpPushDup
	OpPushConst
	OpPushLocal
	OpPushLocal0
	OpPushLocal1
	OpPushGlobal
	OpPushCaptured
	OpPushCollElem

	OpPop
	OpPopLocal
	OpPopLocal0
	OpPopLocal1
	OpPopGlobal
	OpPopCaptured
	OpPopCollElem

	OpBinaryOp
	OpUnaryOp

	OpBranch
	OpBranchIfFalse

	OpCallFunc
	OpCallBuiltin
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
	OpPushConst:    {"PUSH.C", []OperandInfo{{2, OperandConst}}},
	OpPushLocal:    {"PUSH.L", []OperandInfo{{1, OperandLiteral}}},
	OpPushLocal0:   {"PUSH.L0", []OperandInfo{}},
	OpPushLocal1:   {"PUSH.L1", []OperandInfo{}},
	OpPushGlobal:   {"PUSH.G", []OperandInfo{{1, OperandLiteral}}},
	OpPushCaptured: {"PUSH.O", []OperandInfo{{1, OperandLiteral}}},
	OpPushCollElem: {"PUSH.CE", []OperandInfo{}},

	OpPop:         {"POP", []OperandInfo{}},
	OpPopLocal:    {"POP.L", []OperandInfo{{1, OperandLiteral}}},
	OpPopLocal0:   {"POP.L0", []OperandInfo{}},
	OpPopLocal1:   {"POP.L1", []OperandInfo{}},
	OpPopGlobal:   {"POP.G", []OperandInfo{{1, OperandLiteral}}},
	OpPopCaptured: {"POP.O", []OperandInfo{{1, OperandLiteral}}},
	OpPopCollElem: {"POP.CE", []OperandInfo{}},

	OpBinaryOp: {"OP.I", []OperandInfo{{1, OperandBinaryOp}}},
	OpUnaryOp:  {"OP.P", []OperandInfo{{1, OperandUnaryOp}}},

	OpBranch:        {"BR", []OperandInfo{{2, OperandLiteral}}},
	OpBranchIfFalse: {"BR.F", []OperandInfo{{2, OperandLiteral}}},

	OpCallFunc:    {"CALL.F", []OperandInfo{{1, OperandLiteral}}},
	OpCallBuiltin: {"CALL.B", []OperandInfo{{1, OperandBuiltin}, {1, OperandLiteral}}},
	OpReturn:      {"RET", []OperandInfo{}},

	OpNewArray: {"NEW.A", []OperandInfo{{2, OperandLiteral}}},
	OpNewFunc:  {"NEW.F", []OperandInfo{{2, OperandFunc}}},
}
