package application

func (c *CPU) ClearCarryFlag() {
	c.carry = false
}

func (c *CPU) SetCarryFlag() {
	c.carry = true
}

func (c *CPU) ClearInterruptFlag() {
	c.irq = true
}

func (c *CPU) SetInterruptFlag() {
	c.irq = false
}

func (c *CPU) ClearOverflowFlag() {
	c.overflow = false
}

func (c *CPU) SetOverflowFlag() {
	c.overflow = true
}

func (c *CPU) NoOp() {}

func (c *CPU) ClearDecimalFlag() {
	c.decimal = false
}

func (c *CPU) SetDecimalFlag() {
	c.decimal = true
}
