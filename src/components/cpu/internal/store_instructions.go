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
	c.memory.Write(c.GetAddressByIndexedAbsoluteMode(address, register), c.acc)
}

func (c *cpu) StoreAccumulatorIntoMemoryWithIndexZeroPage(address uint8, register string) {
	c.memory.Write(c.GetAddressByIndexedZeroPageMode(address, register), c.acc)
}

func (c *cpu) StoreRegisterIntoMemoryWithIndirectAddress(initialAddress uint16, register string) {
	c.storeRegisterIntoMemory(c.GetAddressByIndirectAbsoluteMode(initialAddress), register)
}

func (c *cpu) StoreRegisterIntoMemoryWithIndexedXIndirectAddress(initialAddress uint8, register string) {
	c.storeRegisterIntoMemory(c.GetAddressByIndexedIndirectXMode(initialAddress), register)
}

func (c *cpu) StoreRegisterIntoMemoryWithIndirectIndexedYAddress(initialAddress uint8, register string) {
	c.storeRegisterIntoMemory(c.GetAddressByIndirectIndexedYMode(initialAddress), register)
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
