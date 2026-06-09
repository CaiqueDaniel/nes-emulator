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
}

func (c *cpu) LoadValueFromMemoryIntoRegister(address uint16, register string) {
	value := c.memory.Read(address)
	c.LoadValueIntoRegister(value, register)
}

func (c *cpu) LoadValueFromMemoryIntoRegisterZeroPage(address uint8, register string) {
	c.LoadValueFromMemoryIntoRegister(uint16(address), register)
}

func (c *cpu) LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(address uint16, register string) {
	value := c.memory.Read(c.GetAddressByIndexedAbsoluteMode(address, register))
	c.acc = value
}

func (c *cpu) LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(address uint8, register string) {
	value := c.memory.Read(c.GetAddressByIndexedZeroPageMode(address, register))
	c.acc = value
}

func (c *cpu) LoadValueFromMemoryIntoRegisterWithIndexedXIndirectAddress(initialAddress uint8, register string) {
	c.LoadValueIntoRegister(c.memory.Read(c.GetAddressByIndexedIndirectXMode(initialAddress)), register)
}

func (c *cpu) LoadValueFromMemoryIntoRegisterWithIndirectIndexedYAddress(initialAddress uint8, register string) {
	c.LoadValueIntoRegister(c.memory.Read(c.GetAddressByIndirectIndexedYMode(initialAddress)), register)
}
