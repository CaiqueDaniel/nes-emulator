package internal

func (c *cpu) StoreRegisterIntoAbsoluteMemory(value uint16, register string) {
	switch register {
	case REGISTER_X:
		c.memory.Write(value, c.x)
	case REGISTER_Y:
		c.memory.Write(value, c.y)
	case ACCUMULATOR:
		c.memory.Write(value, c.acc)
	}
}

func (c *cpu) StoreRegisterIntoZeroPageMemory(value uint8, register string) {
	c.StoreRegisterIntoAbsoluteMemory(uint16(value), register)
}

func (c *cpu) StoreAccumulatorIntoMemoryWithIndexAbsolute(address uint16, register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + uint16(registerValue)
	c.memory.Write(finalAddress, c.acc)
}

func (c *cpu) StoreAccumulatorIntoMemoryWithIndexZeroPage(address uint8, register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + registerValue
	c.memory.Write(uint16(finalAddress), c.acc)
}
