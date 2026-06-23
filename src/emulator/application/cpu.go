package application

type CPU interface {
	AddWithCarry(value uint8)
	And(value uint8)
	ArithmeticShiftLeft()
	ArithmeticShiftLeftAbsolute(address uint16, xIndexedMode bool)
	ArithmeticShiftLeftZeroPage(address uint8, xIndexedMode bool)
	Bit(value uint8)
	BranchIfCarryIsClear(value uint8)
	BranchIfCarryIsSet(value uint8)
	BranchIfEqual(value uint8)
	BranchIfNegative(value uint8)
	BranchIfNotEqual(value uint8)
	BranchIfOverflowClear(value uint8)
	BranchIfOverflowSet(value uint8)
	BranchIfPositive(value uint8)
	Break()
	ClearCarryFlag()
	ClearDecimalFlag()
	ClearInterruptFlag()
	ClearOverflowFlag()
	CompareWithRegister(value uint8, register string)
	DecrementMemory(address uint16)
	DecrementRegister(register string)
	GetCarryFlag() bool
	GetDebugData() map[string]uint8
	GetDecimalFlag() bool
	GetIRQFlag() bool
	GetNegativeFlag() bool
	GetOverflowFlag() bool
	GetProgramCounter() uint16
	GetStackPointer() uint8
	GetValueByAbsoluteMode(address uint16) uint8
	GetValueByIndexedAbsoluteMode(address uint16, index uint8) uint8
	GetValueByIndexedIndirectXMode(initialAddress uint8) uint8
	GetValueByIndirectAbsoluteMode(initialAddress uint16) uint8
	GetValueByIndirectIndexedYMode(initialAddress uint8) uint8
	GetValueByZeroPageIndexedMode(address uint8, index uint8) uint8
	GetValueByZeroPageMode(address uint8) uint8
	GetZeroFlag() bool
	IncrementMemory(address uint16)
	IncrementRegister(register string)
	JumpProgramCounterByIndirectValue(address uint16)
	JumpProgramCounterToSubRoutine(value uint16)
	JumpProgramCounterToValue(value uint16)
	LoadValueIntoRegister(value uint8, register string)
	LogicalShiftRight()
	LogicalShiftRightAbsolute(address uint16, xIndexedMode bool)
	LogicalShiftRightZeroPage(address uint8, xIndexedMode bool)
	NoOp()
	Or(value uint8)
	PullAccFromStack()
	PullStatusFromStack()
	PullValueFromStack() uint8
	PushAccToStack()
	PushFlagsIntoStack()
	PushStatusIntoStack()
	PushValueToStack(value uint8)
	Reset()
	ReturnFromInterrupt()
	ReturnFromSubRoutine()
	RotateLeft()
	RotateLeftAbsolute(address uint16, xIndexedMode bool)
	RotateLeftZeroPage(address uint8, xIndexedMode bool)
	RotateRight()
	RotateRightAbsolute(address uint16, xIndexedMode bool)
	RotateRightZeroPage(address uint8, xIndexedMode bool)
	SetCarryFlag()
	SetDecimalFlag()
	SetInterruptFlag()
	SetOverflowFlag()
	StoreRegisterIntoAbsoluteMemory(value uint16, register string)
	SubtractWithCarry(value uint8)
	TransferFromAccumulatorToRegister(register string)
	TransferFromRegisterToAccumulator(register string)
	TransferStackToX()
	TransferXToStack()
	Xor(value uint8)
	GetTotalCycles() uint64
}
