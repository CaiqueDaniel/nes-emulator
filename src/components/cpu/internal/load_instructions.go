package internal

func (c *cpu) LoadValueIntoRegister(value uint8, register string) {
	switch register {
	case REGISTER_X:
		c.x = value
	case REGISTER_Y:
		c.y = value
	case ACCUMULATOR:
		c.acc = value
	}

	c.zero = c.isValueZero(value)
	c.negative = c.isValueNegative(value)
}
