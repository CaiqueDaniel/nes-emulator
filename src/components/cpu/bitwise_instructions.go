package cpu

func (c *cpu) And(value uint8) {
	c.acc = c.acc & value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Or(value uint8) {
	c.acc = c.acc | value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Xor(value uint8) {
	c.acc = c.acc ^ value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Bit(value uint8) {
	result := c.acc & value

	c.zero = c.isValueZero(result)
	c.negative = c.isValueNegative(result)
	c.overflow = result&0b01000000 != 0
}
