package application

func (c *CPU) GetDebugData() map[string]uint8 {
	return map[string]uint8{
		"acc": c.acc,
		"x":   c.x,
		"y":   c.y,
	}
}

func (c *CPU) GetCarryFlag() bool {
	return c.carry
}

func (c *CPU) GetOverflowFlag() bool {
	return c.overflow
}

func (c *CPU) GetIRQFlag() bool {
	return c.irq
}

func (c *CPU) GetZeroFlag() bool {
	return c.zero
}

func (c *CPU) GetNegativeFlag() bool {
	return c.negative
}

func (c *CPU) GetProgramCounter() uint16 {
	return c.programCounter
}

func (c *CPU) GetDecimalFlag() bool {
	return c.decimal
}

func (c *CPU) GetStackPointer() uint8 {
	return c.stackPointer
}

func (c *CPU) GetNumberOfInstructions() int {
	count := 0

	for _, instruction := range c.instructionSet {
		if instruction == nil {
			continue
		}

		count++
	}

	return count
}

func (c *CPU) GetNMIFlag() bool {
	return c.nmi
}
