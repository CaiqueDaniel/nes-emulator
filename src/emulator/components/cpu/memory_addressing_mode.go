package cpu

func (c *cpu) GetValueByAbsoluteMode(address uint16) uint8 {
	return c.memory.Read(address)
}

func (c *cpu) GetValueByZeroPageMode(address uint8) uint8 {
	return c.memory.Read(uint16(address))
}

func (c *cpu) GetValueByIndexedAbsoluteMode(address uint16, index uint8) uint8 {
	return c.memory.Read(c.GetAddressByIndexedAbsoluteMode(address, index))
}

func (c *cpu) GetValueByZeroPageIndexedMode(address uint8, index uint8) uint8 {
	return c.memory.Read(uint16(address) + uint16(index))
}

func (c *cpu) GetValueByZeroPageIndexedModeWithDummyRead(address uint8, index uint8) uint8 {
	c.doDummyMemoryRead(uint16(address))
	return c.memory.Read(uint16(address) + uint16(index))
}

func (c *cpu) GetValueByIndirectAbsoluteMode(initialAddress uint16) uint8 {
	lastByte := c.memory.Read(initialAddress)
	firstByte := c.memory.Read(initialAddress + 1)

	return c.memory.Read(uint16(firstByte)<<8 | uint16(lastByte))
}

func (c *cpu) GetValueByIndexedIndirectXMode(initialAddress uint8) uint8 {
	c.memory.Read(uint16(initialAddress))
	return c.memory.Read(c.GetAddressByIndexedIndirectXMode(uint16(initialAddress)))
}

func (c *cpu) GetValueByIndirectIndexedYMode(initialAddress uint8) uint8 {
	return c.memory.Read(c.GetAddressByIndirectIndexedYMode(initialAddress))
}

func (c *cpu) GetAddressByIndexedAbsoluteMode(address uint16, index uint8) uint16 {
	prevHighAddress := uint8(address >> 8)
	newAddress := address + uint16(index)
	currentHighAddress := uint8(newAddress >> 8)

	if prevHighAddress != currentHighAddress {
		c.doDummyMemoryRead(address)
	}

	return newAddress
}

func (c *cpu) GetAddressByIndexedAbsoluteModeWithDummyRead(address uint16, index uint8) uint16 {
	c.doDummyMemoryRead(address)
	return c.GetAddressByIndexedAbsoluteMode(address, index)
}

func (c *cpu) GetAddressByIndexedIndirectXModeWithDummyRead(initialAddress uint16) uint16 {
	c.doDummyMemoryRead(initialAddress)
	return c.GetAddressByIndexedIndirectXMode(initialAddress)
}

func (c *cpu) GetAddressByIndexedIndirectXMode(initialAddress uint16) uint16 {
	pivotAddress := uint16(initialAddress) + uint16(c.x)
	lastByte := c.memory.Read(pivotAddress)
	firstByte := c.memory.Read(pivotAddress + 1)

	return uint16(firstByte)<<8 | uint16(lastByte)
}

func (c *cpu) GetAddressByIndirectIndexedYMode(initialAddress uint8) uint16 {
	lowByte := c.memory.Read(uint16(initialAddress))
	highByte := c.memory.Read(uint16(initialAddress) + 1)
	baseAddress := uint16(highByte)<<8 | uint16(lowByte)
	newAddresss := baseAddress + uint16(c.y)

	if uint8(newAddresss>>8) != highByte {
		c.doDummyMemoryRead(baseAddress)
	}

	return newAddresss
}

func (c *cpu) GetAddressByIndirectIndexedYModeWithDummyRead(initialAddress uint8) uint16 {
	lastByte := c.memory.Read(uint16(initialAddress))
	firstByte := c.memory.Read(uint16(initialAddress) + 1)
	incompleteAddress := uint16(firstByte)<<8 | uint16(lastByte)

	c.doDummyMemoryRead(incompleteAddress)

	return incompleteAddress + uint16(c.y)
}
