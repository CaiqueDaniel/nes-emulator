package cpu

import (
	"nes-emu/src/emulator/application"
)

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"
const STACK_START = 0x01FF
const STACK_END = 0x0100

type cpu struct {
	programCounter                                       uint16
	acc, x, y                                            uint8
	carry, zero, overflow, negative, irq, decimal, bFlag bool
	stackPointer                                         uint8
	memory                                               application.Memory
}

func NewCpu(memory application.Memory) application.CPU {
	c := &cpu{
		memory: memory,
	}
	c.Reset()
	return c
}

func NewCpuWithProgramCounter(memory application.Memory, programCounter uint16) application.CPU {
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
