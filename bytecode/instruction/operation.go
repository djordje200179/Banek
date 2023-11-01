package instruction

type Operation byte

const (
	Invalid Operation = iota

	PushDuplicate
	PushConst
	PushLocal
	PushGlobal
	PushCaptured
	PushBuiltin
	PushCollectionElement

	Pop
	PopLocal
	PopGlobal
	PopCaptured
	PopCollectionElement

	OperationInfix
	OperationPrefix

	Branch
	BranchIfFalse

	Call
	Return

	NewArray
	NewFunction
)

func (operation Operation) String() string {
	return operation.Info().Name
}

func (operation Operation) Info() OperationInfo {
	if operation < 0 || operation >= Operation(len(operationInfos)) {
		return operationInfos[Invalid]
	}

	return operationInfos[operation]
}

type OperationInfo struct {
	Name     string
	Operands []OperandInfo
}

func (opInfo OperationInfo) Size() int {
	size := 1

	for _, operand := range opInfo.Operands {
		size += operand.Width
	}

	return size
}

func (opInfo OperationInfo) OperandOffset(index int) int {
	offset := 1
	for i := 0; i < index; i++ {
		offset += opInfo.Operands[i].Width
	}

	return offset
}

var operationInfos = []OperationInfo{
	Invalid: {"INVALID", []OperandInfo{}},

	PushDuplicate:         {"PUSH.D", []OperandInfo{}},
	PushConst:             {"PUSH.C", []OperandInfo{{2, Constant}}},
	PushLocal:             {"PUSH.L", []OperandInfo{{1, Literal}}},
	PushGlobal:            {"PUSH.G", []OperandInfo{{1, Literal}}},
	PushCaptured:          {"PUSH.O", []OperandInfo{{1, Literal}}},
	PushBuiltin:           {"PUSH.B", []OperandInfo{{1, Literal}}},
	PushCollectionElement: {"PUSH.CE", []OperandInfo{}},

	Pop:                  {"POP", []OperandInfo{}},
	PopLocal:             {"POP.L", []OperandInfo{{1, Literal}}},
	PopGlobal:            {"POP.G", []OperandInfo{{1, Literal}}},
	PopCaptured:          {"POP.O", []OperandInfo{{1, Literal}}},
	PopCollectionElement: {"POP.CE", []OperandInfo{}},

	OperationInfix:  {"OP.I", []OperandInfo{{1, Literal}}},
	OperationPrefix: {"OP.P", []OperandInfo{{1, Literal}}},

	Branch:        {"BR", []OperandInfo{{2, Literal}}},
	BranchIfFalse: {"BR.F", []OperandInfo{{2, Literal}}},

	Call:   {"CALL", []OperandInfo{{1, Literal}}},
	Return: {"RET", []OperandInfo{}},

	NewArray:    {"NEW.A", []OperandInfo{{2, Literal}}},
	NewFunction: {"NEW.F", []OperandInfo{{2, Function}}},
}
