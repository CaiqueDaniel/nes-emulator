package internal

import "nes-emu/src/components/memory"

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"

type cpu struct {
	programCounter                                  uint16
	acc, x, y                                       uint8
	carry, zero, overflow, negative, interrupt, irq bool
	stackPointer                                    uint8
	memory                                          memory.Memory
}

func NewCpu(memory memory.Memory) *cpu {
	c := &cpu{
		memory: memory,
	}
	c.Reset()
	return c
}

func NewCpuWithProgramCounter(memory memory.Memory, programCounter uint16) *cpu {
	c := &cpu{
		memory:         memory,
		programCounter: programCounter,
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

func (c *cpu) GetInterruptFlag() bool {
	return c.interrupt
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
