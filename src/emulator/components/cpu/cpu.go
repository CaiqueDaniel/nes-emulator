package cpu

import (
	"nes-emu/src/emulator/application"
	"slices"
)

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"
const STACK_START = 0x01FF
const STACK_END = 0x0100
const START_POINTER = 0xFFFC

type cpu struct {
	programCounter                                            uint16
	acc, x, y                                                 uint8
	carry, zero, overflow, negative, irq, decimal, bFlag, nmi bool
	stackPointer                                              uint8
	memory                                                    application.Memory
	bus                                                       application.Bus
	stopPcAt                                                  int
	instructionSet                                            instructionSet
}

func NewCpu(memory application.Memory, bus application.Bus) application.CPU {
	c := &cpu{
		memory:   memory,
		bus:      bus,
		stopPcAt: -1,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func NewCpuWithStopAt(memory application.Memory, bus application.Bus, stopPcAt int) *cpu {
	c := &cpu{
		memory:   memory,
		bus:      bus,
		stopPcAt: stopPcAt,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func NewCpuWithProgramCounter(memory application.Memory, programCounter uint16, bus application.Bus) *cpu {
	c := &cpu{
		memory:         memory,
		programCounter: programCounter,
		stackPointer:   0xFF,
		bus:            bus,
		stopPcAt:       -1,
	}
	c.instructionSet = c.initInstructions()
	return c
}

func NewCpuWithInternal(memory application.Memory, bus application.Bus) *cpu {
	c := &cpu{
		memory:   memory,
		bus:      bus,
		stopPcAt: -1,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func (c *cpu) RunProgram() {
	startAddressLow := c.readFromMemory(START_POINTER)
	startAddressHigh := c.readFromMemory(START_POINTER + 1)
	startAddress := (uint16(startAddressHigh) << 8) + uint16(startAddressLow)
	c.programCounter = startAddress

	for {
		if c.nmi {
			c.HandleNMI()
			c.nmi = false
		}

		opCode := c.readFromMemory(c.programCounter)
		c.programCounter++

		c.interpretInstruction(opCode)

		if c.stopPcAt != -1 && c.programCounter >= uint16(c.stopPcAt) {
			break
		}
	}
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
	c.irq = true
	c.stackPointer = 0xFF
	c.nmi = false
}

func (c *cpu) SetNMI() {
	c.nmi = true
}

func (c *cpu) HandleNMI() {
	const nmi_handler_byte_ptr = 0xFFFA

	c.PushStatusIntoStack()

	lowByte := c.readFromMemory(nmi_handler_byte_ptr)
	highByte := c.readFromMemory(nmi_handler_byte_ptr + 1)

	c.programCounter = uint16(highByte)<<8 | uint16(lowByte)
}

func (c *cpu) PushValueToStack(value uint8) {
	c.writeToMemory(STACK_END+uint16(c.stackPointer), value)
	c.stackPointer--
}

func (c *cpu) PullValueFromStack() uint8 {
	c.stackPointer++
	value := c.readFromMemory(STACK_END + uint16(c.stackPointer))

	c.writeToMemory(STACK_END+uint16(c.stackPointer), 0)

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

func (c *cpu) interpretInstruction(opCode uint8) {
	instruction := c.instructionSet[opCode]

	if instruction == nil {
		return
	}

	bytes := make([]uint8, 0)

	for i := 0; i < instruction.ArgsBytes; i++ {
		bytes = append(bytes, c.readFromMemory(c.programCounter))
		c.programCounter++
	}

	if instruction.ArgsBytes == 0 {
		c.doDummyMemoryRead(c.programCounter)
	}

	slices.Reverse(bytes)
	instruction.Method(bytes)
}

func (c *cpu) writeToMemory(address uint16, value uint8) {
	c.memory.Write(address, value)
}

func (c *cpu) readFromMemory(address uint16) uint8 {
	return c.memory.Read(address)
}

func (c *cpu) doDummyMemoryRead(address uint16) {
	c.readFromMemory(address)
}

func transformFlagIntoUint8(flag bool) uint8 {
	if flag {
		return 1
	}
	return 0
}
