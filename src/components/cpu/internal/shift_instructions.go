package internal

func (c *cpu) ArithmeticShiftLeft() {
	prevValue := c.acc
	c.acc = c.acc << 1
	c.negative = c.isValueNegative(c.acc)
	c.carry = c.acc < prevValue
}

func (c *cpu) ArithmeticShiftLeftAbsolute(address uint16, xIndexedMode bool) {
	if xIndexedMode {
		address = c.getAddressByIndexedAbsoluteMode(address, c.x)
	}

	prevValue := c.GetValueByAbsoluteMode(address)
	value := prevValue << 1

	c.memory.Write(address, prevValue)
	c.memory.Write(address, value)

	c.negative = c.isValueNegative(value)
	c.carry = value < prevValue
}

func (c *cpu) ArithmeticShiftLeftZeroPage(address uint8, xIndexedMode bool) {
	c.ArithmeticShiftLeftAbsolute(uint16(address), xIndexedMode)
}
