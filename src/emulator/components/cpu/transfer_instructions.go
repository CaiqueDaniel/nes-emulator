package cpu

func (c *cpu) TransferFromAccumulatorToRegister(register string) {
	switch register {
	case REGISTER_X:
		c.x = c.acc
	case REGISTER_Y:
		c.y = c.acc
	}

	c.negative = c.isValueNegative(c.acc)
	c.zero = c.isValueZero(c.acc)
}

func (c *cpu) TransferFromRegisterToAccumulator(register string) {
	switch register {
	case REGISTER_X:
		c.acc = c.x
	case REGISTER_Y:
		c.acc = c.y
	}

	c.negative = c.isValueNegative(c.acc)
	c.zero = c.isValueZero(c.acc)
}
