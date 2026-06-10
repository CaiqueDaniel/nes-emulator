package internal

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
	c.irq = true
}
