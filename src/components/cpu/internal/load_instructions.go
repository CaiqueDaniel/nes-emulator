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
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + uint16(registerValue)
	value := c.memory.Read(finalAddress)
	c.acc = value
}

func (c *cpu) LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(address uint8, register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + registerValue
	value := c.memory.Read(uint16(finalAddress))
	c.acc = value
}

func (c *cpu) LoadValueFromMemoryIntoRegisterWithIndexedXIndirectAddress(initialAddress uint8, register string) {
	pivotAddress := uint16(initialAddress) + uint16(c.x)
	lastByte := c.memory.Read(pivotAddress)
	firstByte := c.memory.Read(pivotAddress + 1)
	finalAddress := uint16(firstByte)<<8 | uint16(lastByte)

	c.LoadValueIntoRegister(c.memory.Read(finalAddress), register)
}

func (c *cpu) LoadValueFromMemoryIntoRegisterWithIndirectIndexedYAddress(initialAddress uint8, register string) {
	lastByte := c.memory.Read(uint16(initialAddress))
	firstByte := c.memory.Read(uint16(initialAddress) + 1)
	finalAddress := uint16(firstByte)<<8 | uint16(lastByte) + uint16(c.y)

	c.LoadValueIntoRegister(c.memory.Read(finalAddress), register)
}
