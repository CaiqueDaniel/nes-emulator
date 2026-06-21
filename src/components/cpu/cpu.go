package cpu

import "nes-emu/src/components/memory"

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"
const STACK_START = 0x01FF
const STACK_END = 0x0100

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
}

type cpu struct {
	programCounter                                       uint16
	acc, x, y                                            uint8
	carry, zero, overflow, negative, irq, decimal, bFlag bool
	stackPointer                                         uint8
	memory                                               memory.Memory
}

func NewCpu(memory memory.Memory) CPU {
	c := &cpu{
		memory: memory,
	}
	c.Reset()
	return c
}

func NewCpuWithProgramCounter(memory memory.Memory, programCounter uint16) CPU {
	c := &cpu{
		memory:         memory,
		programCounter: programCounter,
		stackPointer:   0xFF,
	}
	return c
}

func (c *cpu) Reset() {
	c.programCounter = 0
	c.acc = 0
	c.x = 0
	c.y = 0
	c.carry = false
	c.zero = false
	c.overflow = false
	c.negative = false
	c.stackPointer = 0xFF
}

func (c *cpu) PushValueToStack(value uint8) {
	c.memory.Write(STACK_END+uint16(c.stackPointer), value)
	c.stackPointer--
}

func (c *cpu) PullValueFromStack() uint8 {
	c.stackPointer++
	value := c.memory.Read(STACK_END + uint16(c.stackPointer))

	c.memory.Write(STACK_END+uint16(c.stackPointer), 0)

	return value
}

func (c *cpu) PushFlagsIntoStack() {
	carry := transformFlagIntoUint8(c.carry)
	zero := transformFlagIntoUint8(c.zero) << 1
	irq := transformFlagIntoUint8(c.irq) << 2
	decimal := transformFlagIntoUint8(c.decimal) << 3
	breakFlag := transformFlagIntoUint8(c.bFlag) << 4
	overflow := transformFlagIntoUint8(c.overflow) << 6
	negative := transformFlagIntoUint8(c.negative) << 7

	valueToPush := 0b00100000 | carry | zero | irq | decimal | breakFlag | overflow | negative
	c.PushValueToStack(valueToPush)
}

func (c *cpu) GetDebugData() map[string]uint8 {
	return map[string]uint8{
		"acc": c.acc,
		"x":   c.x,
		"y":   c.y,
	}
}

func (c *cpu) GetCarryFlag() bool {
	return c.carry
}

func (c *cpu) GetOverflowFlag() bool {
	return c.overflow
}

func (c *cpu) GetIRQFlag() bool {
	return c.irq
}

func (c *cpu) GetZeroFlag() bool {
	return c.zero
}

func (c *cpu) GetNegativeFlag() bool {
	return c.negative
}

func (c *cpu) GetProgramCounter() uint16 {
	return c.programCounter
}

func (c *cpu) GetDecimalFlag() bool {
	return c.decimal
}

func transformFlagIntoUint8(flag bool) uint8 {
	if flag {
		return 1
	}
	return 0
}

func (c *cpu) GetStackPointer() uint8 {
	return c.stackPointer
}

//Arithmetic		SBC
