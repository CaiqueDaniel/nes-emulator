package cpu

func (c *cpu) StoreRegisterIntoAbsoluteMemory(value uint16, register string) {
	switch register {
	case REGISTER_X:
		c.writeToMemory(value, c.x)
	case REGISTER_Y:
		c.writeToMemory(value, c.y)
	case ACCUMULATOR:
		c.writeToMemory(value, c.acc)
	}
}

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
