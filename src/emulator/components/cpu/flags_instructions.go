package cpu

func (c *cpu) ClearCarryFlag() {
	c.carry = false
	c.tick(1)
}

func (c *cpu) SetCarryFlag() {
	c.carry = true
	c.tick(1)
}

func (c *cpu) ClearInterruptFlag() {
	c.irq = true
	c.tick(1)
}

func (c *cpu) SetInterruptFlag() {
	c.irq = false
	c.tick(1)
}

func (c *cpu) ClearOverflowFlag() {
	c.overflow = false
	c.tick(1)
}

func (c *cpu) SetOverflowFlag() {
	c.overflow = true
	c.tick(1)
}

func (c *cpu) NoOp() {
	c.tick(1)
}

func (c *cpu) Break() {
	highAddress := uint8(c.programCounter >> 8)
	lowAddress := uint8(c.programCounter)

	c.irq = true
	c.bFlag = true
	c.PushValueToStack(highAddress)
	c.PushValueToStack(lowAddress)
	c.PushStatusIntoStack()

	c.programCounter = uint16(c.memory.Read(0xFFFE))
	//c.tick(7)
}

func (c *cpu) ClearDecimalFlag() {
	c.decimal = false
	c.tick(1)
}

func (c *cpu) SetDecimalFlag() {
	c.decimal = true
	c.tick(1)
}
