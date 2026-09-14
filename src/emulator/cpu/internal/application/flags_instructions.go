package application

func (c *CPU) ClearCarryFlag() {
	c.carry = false
}

func (c *CPU) SetCarryFlag() {
	c.carry = true
}

func (c *CPU) ClearInterruptFlag() {
	c.interrupt = false
}

func (c *CPU) SetInterruptFlag() {
	c.interrupt = true
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
