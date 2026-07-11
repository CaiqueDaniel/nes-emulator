package cpu

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

func (c *cpu) GetTotalCycles() uint64 {
	return c.totalCycles
}

func (c *cpu) GetStackPointer() uint8 {
	return c.stackPointer
}

func (c *cpu) GetNumberOfInstructions() int {
	count := 0

	for _, instruction := range c.instructionSet {
		if instruction == nil {
			continue
		}

		count++
	}

	return count
}
