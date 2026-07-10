package cpu

func (c *cpu) ClearCarryFlag() {
	c.carry = false
}

func (c *cpu) SetCarryFlag() {
	c.carry = true
}

func (c *cpu) ClearInterruptFlag() {
	c.irq = true
}

func (c *cpu) SetInterruptFlag() {
	c.irq = false
}

func (c *cpu) ClearOverflowFlag() {
	c.overflow = false
}

func (c *cpu) SetOverflowFlag() {
	c.overflow = true
}

func (c *cpu) NoOp() {}

func (c *cpu) Break() {
	highAddress := uint8(c.programCounter >> 8)
	lowAddress := uint8(c.programCounter)

	c.tick(1)

	c.irq = true
	c.bFlag = true

	c.tick(1)

	c.PushValueToStack(highAddress)
	c.PushValueToStack(lowAddress)
	c.PushStatusIntoStack()

	c.programCounter = uint16(c.memory.Read(0xFFFE))
}

func (c *cpu) ClearDecimalFlag() {
	c.decimal = false
}

func (c *cpu) SetDecimalFlag() {
	c.decimal = true
}
