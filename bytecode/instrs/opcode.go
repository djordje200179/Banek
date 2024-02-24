package instrs

type Opcode byte

const (
	OpHalt Opcode = iota
	OpCall
	OpReturn
)

const (
	OpJump Opcode = 0x04 + iota
	OpBranchFalse
	OpBranchTrue
)

const (
	OpLoadGlobal Opcode = 0x08 + iota
	OpLoadCaptured
	OpLoadLocal
	OpLoadLocal0
	OpLoadLocal1
	OpLoadLocal2
)

const (
	OpConst0 Opcode = 0x10 + iota
	OpConst1
	OpConst2
	OpConst3
	OpConstN1
	OpConstInt
	OpConstStr
	OpConstTrue
	OpConstFalse
	OpConstUndef
	OpBuiltin
)

const (
	OpMakeArray Opcode = 0x1C + iota
	OpNewArray
	OpMakeFunc
)

const (
	OpPop Opcode = 0x20 + iota
	OpStoreGlobal
	OpStoreCaptured
	OpStoreLocal
	OpStoreLocal0
	OpStoreLocal1
	OpStoreLocal2
)

const (
	OpCollGet Opcode = 0x28 + iota
	OpCollSet
)

const (
	OpDup Opcode = 0x2C + iota
	OpDup2
	OpSwap
)

const (
	OpCompEq Opcode = 0x30 + iota
	OpCompNe
	OpCompLt
	OpCompLe
	OpCompGt
	OpCompGe
)

const (
	OpNeg Opcode = 0x38 + iota
	OpNot
)

const (
	OpAdd Opcode = 0x40 + iota
	OpSub
	OpMul
	OpDiv
	OpMod
)

func (opcode Opcode) String() string {
	return opcode.Info().Name
}
func (opcode Opcode) Info() InstrInfo {
	return instrInfos[opcode]
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
	for i := range index {
		offset += instrInfo.Operands[i].Width
	}

	return offset
}
