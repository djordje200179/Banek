package instruction

type Operation byte

const (
	Invalid Operation = iota

	PushConst
	Pop

	Negate
	Add
	Subtract
	Multiply
	Divide

	Equals
	NotEquals
	LessThan
	LessThanOrEquals

	Branch
	BranchIfFalse

	NewArray
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

	PushConst: {"PUSH.C", []OperandInfo{constantPoolOperand}},
	Pop:       {"POP", []OperandInfo{}},

	Negate:   {"NEG", []OperandInfo{}},
	Add:      {"ADD", []OperandInfo{}},
	Subtract: {"SUB", []OperandInfo{}},
	Multiply: {"MUL", []OperandInfo{}},
	Divide:   {"DIV", []OperandInfo{}},

	Equals:           {"EQ", []OperandInfo{}},
	NotEquals:        {"NEQ", []OperandInfo{}},
	LessThan:         {"LT", []OperandInfo{}},
	LessThanOrEquals: {"LTE", []OperandInfo{}},

	Branch:        {"BR", []OperandInfo{{2, Literal}}},
	BranchIfFalse: {"BR.F", []OperandInfo{{2, Literal}}},

	NewArray: {"NEW.A", []OperandInfo{{2, Literal}}},
}
