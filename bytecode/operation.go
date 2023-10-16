package bytecode

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
	GreaterThan
	GreaterThanOrEquals
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
	Name          string
	OperandWidths []int
}

func (opInfo OperationInfo) OperandsSize() int {
	size := 0
	for _, width := range opInfo.OperandWidths {
		size += width
	}

	return size
}

var operationInfos = []OperationInfo{
	Invalid: {"INVALID", []int{}},

	PushConst: {"PUSH.C", []int{2}},
	Pop:       {"POP", []int{}},

	Negate:   {"NEG", []int{}},
	Add:      {"ADD", []int{}},
	Subtract: {"SUB", []int{}},
	Multiply: {"MUL", []int{}},
	Divide:   {"DIV", []int{}},

	Equals:              {"EQ", []int{}},
	NotEquals:           {"NEQ", []int{}},
	LessThan:            {"LT", []int{}},
	LessThanOrEquals:    {"LTE", []int{}},
	GreaterThan:         {"GT", []int{}},
	GreaterThanOrEquals: {"GTE", []int{}},
}
