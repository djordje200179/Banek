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
	OpPushBuiltin Opcode = 0x08 + iota
	OpPushGlobal
	OpPushCaptured
	OpPushCollElem
	OpPushLocal
	OpPushLocal0
	OpPushLocal1
	OpPushLocal2
)

const (
	OpPush0 Opcode = 0x10 + iota
	OpPush1
	OpPush2
	OpPush3
	OpPushN1
	OpPushInt
	OpPushStr
	OpPushTrue
	OpPushFalse
	OpPushUndef
)

const (
	OpPop Opcode = 0x20 + iota
	OpPopGlobal
	OpPopCaptured
	OpPopCollElem
	OpPopLocal
	OpPopLocal0
	OpPopLocal1
	OpPopLocal2
)

const (
	OpDup Opcode = 0x2C + iota
	OpDup2
	OpDup3
	OpSwap
)

const (
	OpBinaryAdd Opcode = 0x30 + iota
	OpBinarySub
	OpBinaryMul
	OpBinaryDiv
	OpBinaryMod
	OpBinaryEq
	OpBinaryNe
	OpBinaryLt
	OpBinaryLe
	OpBinaryGt
	OpBinaryGe
)

const (
	OpUnaryNeg Opcode = 0x3C + iota
	OpUnaryNot
)

const (
	OpMakeArray Opcode = 0x40 + iota
	OpNewArray
	OpMakeFunc
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
	for i := 0; i < index; i++ {
		offset += instrInfo.Operands[i].Width
	}

	return offset
}
