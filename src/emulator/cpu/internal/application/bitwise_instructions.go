package application

func (c *CPU) And(value uint8) {
	c.acc = c.acc & value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *CPU) Or(value uint8) {
	c.acc = c.acc | value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *CPU) Xor(value uint8) {
	c.acc = c.acc ^ value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *CPU) Bit(value uint8) {
	result := c.acc & value

	c.zero = c.isValueZero(result)
	c.negative = c.isValueNegative(result)
	c.overflow = result&0b01000000 != 0
}
