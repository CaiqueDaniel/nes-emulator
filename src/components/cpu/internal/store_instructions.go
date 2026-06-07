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

func (c *cpu) StoreRegisterIntoMemoryWithIndirectAddress(initialAddress uint16, register string) {
	lastByte := c.memory.Read(initialAddress)
	firstByte := c.memory.Read(initialAddress + 1)
	finalAddress := uint16(firstByte)<<8 | uint16(lastByte)

	c.storeRegisterIntoMemory(finalAddress, register)
}

func (c *cpu) StoreRegisterIntoMemoryWithIndexedXIndirectAddress(initialAddress uint8, register string) {
	pivotAddress := uint16(initialAddress) + uint16(c.x)
	lastByte := c.memory.Read(pivotAddress)
	firstByte := c.memory.Read(pivotAddress + 1)
	finalAddress := uint16(firstByte)<<8 | uint16(lastByte)

	c.storeRegisterIntoMemory(finalAddress, register)
}

func (c *cpu) StoreRegisterIntoMemoryWithIndirectIndexedYAddress(initialAddress uint8, register string) {
	lastByte := c.memory.Read(uint16(initialAddress))
	firstByte := c.memory.Read(uint16(initialAddress) + 1)
	finalAddress := uint16(firstByte)<<8 | uint16(lastByte) + uint16(c.y)

	c.storeRegisterIntoMemory(finalAddress, register)
}

func (c *cpu) storeRegisterIntoMemory(address uint16, register string) {
	switch register {
	case REGISTER_X:
		c.memory.Write(address, c.x)
	case REGISTER_Y:
		c.memory.Write(address, c.y)
	case ACCUMULATOR:
		c.memory.Write(address, c.acc)
	}
}
