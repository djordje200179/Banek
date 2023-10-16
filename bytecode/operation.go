package bytecode

type Operation byte

const (
	Invalid Operation = iota

	PushConst

	Add
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
	Invalid: {"Invalid", []int{}},

	PushConst: {"PushConst", []int{2}},

	Add: {"Add", []int{}},
}
